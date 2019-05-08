import os
from datetime import timedelta, datetime

import emails
import json
import logging
import pathlib
from emails.template import JinjaTemplate as T
from prometheus_http_client import Prometheus

# https://prometheus.io/docs/prometheus/latest/querying/basics/#range-vector-selectors
# Issue with leap year
PROM_UNITS = {'s': 1, 'm': 60, 'h': 3600, 'd': 86400, 'w': 604800, 'y': 31536000}
NOT_DEF = '_not_defined'

def main():
    logging.getLogger().setLevel(logging.INFO)
    report_public_envvar()
    data = get_statistics_data_for_template()
    send_email(data)
    return


def report_public_envvar():
    for evar in ('PROMETHEUS_URL', 'REP_TIMEWINDOW', 'REP_SMTP_HOST', 'REP_SMTP_HOST_PORT',
                 'REP_SMTP_USER', 'REP_SMTP_SSL', 'REP_FROM', 'REP_RECIPIENTS'):
        logging.info("%s: %s", evar, os.getenv(evar, "undefined"))


def get_smtp_dict():
    smtp = None
    if os.getenv('REP_SMTP_HOST', None):
        smtp = {
            'host': os.getenv('REP_SMTP_HOST'),
            'port': int(os.getenv('REP_SMTP_HOST_PORT', 25)),
        }

        if os.getenv('REP_SMTP_USER'):
            smtp['user'] = os.getenv('REP_SMTP_USER')
            smtp['password'] = os.getenv('REP_SMTP_USER_PASSWORD')
        smtp['ssl'] = os.getenv('REP_SMTP_SSL', 'xx').lower() in ['true', '1', 't', 'y', 'yes']

    return smtp


def get_recipients():
    return os.getenv('REP_RECIPIENTS').split(',')


def get_from():
    return os.getenv('REP_FROM', 'noreply@death.star.example')


def get_timewindow():
    """

    :return: string
    """
    return os.getenv("REP_TIMEWINDOW", "1d")


def get_timewindow_td():
    """

    :return: datetime.timedelta
    """
    return convert_to_timedelta(get_timewindow())


def send_email(data):
    tpl = None

    with open(pathlib.Path(__file__).parent / 'templates/e-mail-grouped.jinja2') as f:
        tpl = f.read()

    message = emails.Message(html=T(tpl),
                             text="Contains HTML Events summary.",
                             subject=T('Events summary'),
                             mail_from=get_from())

    send_status = message.send(
        to=get_recipients(),
        render={'statistics': data,
                'datetime_from': (datetime.utcnow() - get_timewindow_td()).strftime('%Y-%m-%d %H:%M%z'),
                'datetime_generated': datetime.utcnow().strftime('%Y-%m-%d %H:%M%z'),
                'window': get_timewindow(),
                },
        smtp=get_smtp_dict())

    if send_status.error and get_smtp_dict():
        logging.error(send_status.error)
    else:
        logging.info(message.html_body.replace('  ', ' ').replace('\n', ''))


def get_statistics_data_for_template():
    # Client initializes with env. variable PROMETHEUS_URL
    prom_client = Prometheus()
    timewindow = get_timewindow()
    data_str = prom_client.query(metric="ceil(increase(events_summary_total{}[%s]))" % (timewindow,))

    # Convert data for template
    data = extract_and_convert_data(data_str)
    return data


def extract_and_convert_data(data_str):
    """
    [
        {
            'kind': 'Pod',
            'reason': 'Started',
            'ns': [
                {'name': 'default', 'value': 14},
                {'name': 'kube-system', 'value': 3}
            ]
        },
        {
            'kind': 'Pod',
            'reason': 'Started',
            'ns': [
                {'name': 'kube-system', 'value': 0}
            ]
        },
        {
            'kind': 'ReplicaSet',
            'reason': 'SuccessfulCreate',
            'ns': [
                {'name': 'kube-system', 'value': 0}],
        }
    ]
    """

    data = json.loads(data_str)

    data_iter = ({
        'ns': i['metric'].get('namespace', NOT_DEF),
        'kind': i['metric'].get('kind', NOT_DEF),
        'reason': i['metric'].get('reason', NOT_DEF),
        'value': int(float(i['value'][1])),
    } for i in data['data']['result'])

    data_struct = []
    kind, reason, grouped = None, None, None

    for row in sorted(data_iter, key=lambda t: t["kind"] + t["reason"] + t["ns"]):
        if kind != row['kind'] or reason != row['reason'] or grouped is None:
            grouped = {'kind': row['kind'], 'reason': row['reason'], 'ns': []}
            kind, reason = row['kind'], row['reason']
            data_struct.append(grouped)
        grouped['ns'].append({'name': row['ns'], 'value': row['value']})

    return data_struct


def convert_to_timedelta(s):
    """

    :return: datetime.timedelta
    """
    return timedelta(seconds=int(s[:-1]) * PROM_UNITS[s[-1]])


if __name__ == '__main__':
    main()
