import json
from faker import Faker
import random
from datetime import datetime, timedelta


def check_overlap(new_interval, intervals):
    new_start_time = datetime.strptime(new_interval.split('-')[0], '%H:%M')
    new_end_time = datetime.strptime(new_interval.split('-')[1], '%H:%M')
    if new_start_time == new_end_time or new_start_time > new_end_time:
        return True
    for interval in intervals:
        start_time = datetime.strptime(interval.split('-')[0], '%H:%M')
        end_time = datetime.strptime(interval.split('-')[1], '%H:%M')
        if start_time <= new_start_time < end_time or start_time < new_end_time <= end_time or new_start_time <= start_time < new_end_time or new_start_time < end_time <= new_end_time:
            return True
    return False


def generate_hours(f):
    hours_count = random.randint(1, 4)
    hours = []
    for _ in range(hours_count):
        while True:
            start_time = f.date_time_between(start_date='-1y', end_date='now').strftime('%H:%M')
            end_time = (datetime.strptime(start_time, '%H:%M') + timedelta(
                minutes=random.randint(1, 180))).strftime(
                '%H:%M')
            time_interval = start_time + '-' + end_time
            if not check_overlap(time_interval, hours):
                hours.append(time_interval)
                break
    return hours


def generate_order():
    f = Faker()
    orders = {"orders": []}

    for i in range(100000):
        order = {"cost": f.random_int(min=100, max=10000)}
        delivery_hours = generate_hours(f)
        order["delivery_hours"] = delivery_hours
        order["regions"] = f.random_int(min=1, max=20)
        order["weight"] = round(random.uniform(1, 40), 2)
        orders["orders"].append(order)

    # write to file
    with open("jsons/orders.json", "w") as file:
        json.dump(orders, file, indent=2)
