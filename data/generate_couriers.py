import json
from faker import Faker
import random
from datetime import datetime, timedelta
from generate_orders import check_overlap, generate_hours


def generate_couriers():
    f = Faker()
    couriers = {"couriers": []}
    for i in range(100):
        courier = {"working_hours": generate_hours(f)}
        working_areas_count = random.randint(1, 10)
        working_areas = []
        for _ in range(working_areas_count):
            while True:
                area = f.random_int(min=1, max=10)
                if area not in working_areas:
                    working_areas.append(area)
                    break
        courier["regions"] = working_areas
        types = ["FOOT", "BIKE", "AUTO"]
        courier["courier_type"] = random.choice(types)
        couriers["couriers"].append(courier)

    # write to file
    with open("jsons/couriers.json", "w") as file:
        json.dump(couriers, file, indent=2)
