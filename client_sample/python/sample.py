#!/usr/bin/env python

import requests

data = {
    'Host': '10.10.10.21',
    'Port': 3307,
    'Username': 'root',
    'Password': 'root',
    'Database': 'employees',
    'Sqls': 'alter table employees add column age1 int not null; delete from employees WHERE id = 1;',
}

url = 'http://10.10.10.55:18080/sqlReview'

r = requests.post(url, json = data)

print(r.text)