import os
import json
import uuid
from datetime import datetime
from time import sleep
import random

TABLE_NAME = "artists"

def generate_data():
    """Генерация данных для таблицы artists."""
    genres = ["Pop", "Rock", "Hip-Hop", "Jazz", "Classical", "Electronic"]
    countries = ["USA", "UK", "Canada", "Germany", "France", "Japan", "Australia"]
    
    # Случайное количество записей от 1 до 20
    n = random.randint(1, 20)
    
    return [
        {
            "id": str(uuid.uuid4()),
            "name": f"Artist {random.randint(1, 1000)}",
            "genre": random.choice(genres),
            "country": random.choice(countries),
            "debut_year": random.randint(1950, 2023)
        } 
        for _ in range(n)
    ]

def save_file(data, table_name, file_format="json"):
    """Сохранение данных в файл."""
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    file_id = uuid.uuid4().hex
    filename = f"{file_id}_{table_name}_{timestamp}.{file_format}"

    # Создание директории для файлов
    os.makedirs("input", exist_ok=True)
    filepath = os.path.join("input", filename)

    try:
        if file_format == "json":
            with open(filepath, "w") as f:
                json.dump(data, f, indent=4)
        elif file_format == "csv":
            with open(filepath, "w", newline="") as f:
                fieldnames = ["id", "name", "genre", "country", "debut_year"]
                writer = csv.DictWriter(f, fieldnames=fieldnames)
                writer.writeheader()
                writer.writerows(data)
        else:
            raise ValueError(f"Unsupported file format: {file_format}")
        print(f"File generated: {filepath}")
    except Exception as e:
        print(f"Error while saving file: {e}")
    return filepath

def main():
    """Основная функция."""
    table_name = TABLE_NAME
    file_format = "json"  # Формат файла
    
    while True:
        try:
            print("Generating data...")
            data = generate_data()
            save_file(data, table_name, file_format)
            print("Data saved successfully.")
            sleep(10)  # Интервал 5 минут
        except KeyboardInterrupt:
            print("\nProcess interrupted by user.")
            break
        except Exception as e:
            print(f"Unexpected error: {e}")

if __name__ == "__main__":
    main()
