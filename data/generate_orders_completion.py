import json
from faker import Faker
import random
from datetime import datetime, timedelta


def generate_orders_completion():
    f = Faker()
    orders_completion = {"complete_info": []}
    for i in range(100000):
        order = {"order_id": i + 1,
                 "courier_id": random.randint(1, 100000),
                 "complete_time": f.date_time_between(start_date='-1y', end_date='now').strftime('%Y-%m-%dT%H:%M:%SZ')}
        orders_completion["complete_info"].append(order)

    # write to file
    with open("jsons/orders_completion.json", "w") as file:
        json.dump(orders_completion, file, indent=2)
