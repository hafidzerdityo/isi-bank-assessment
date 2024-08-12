## Consume all of the data

# import json
# import redis, time

# # Hardcoded values
# TOPIC = "jurnalku"
# BROKER = "redis://localhost:6379"
# GROUP = "consumerpython"

# # Initialize the Redis client
# redis_client = redis.Redis.from_url(BROKER)

# # Initialize the last ID to the special ID "0" which means start from the very beginning.
# last_id = "0"

# while True:
#     time.sleep(0.25)
#     try:
#         entries = redis_client.xread({TOPIC: last_id}, count=1, block=0)
#     except redis.RedisError as e:
#         print(f"Failed to read the stream: {str(e)}")
#         continue

#     for entry in entries:
#         message_id, values = entry[1][0]

#         message_json = values[b'message']
#         try:
#             req_redis = json.loads(message_json)
#         except json.JSONDecodeError as e:
#             print(f"Error unmarshalling data: {str(e)}")
#             continue

#         print(f"Message: {req_redis}")

#         # Update last_id to the ID of the last successfully processed message.
#         last_id = message_id


# Consume only the newest data
import json
import redis
import time

# Hardcoded values
TOPIC = "jurnalku"
BROKER = "redis://localhost:6379"
GROUP = "consumerpython"

# Initialize the Redis client
redis_client = redis.Redis.from_url(BROKER)

while True:
    time.sleep(0.25)
    try:
        # Read from the stream, start from the latest message
        entries = redis_client.xread({TOPIC: '$'}, count=1, block=0)
    except redis.RedisError as e:
        print(f"Failed to read the stream: {str(e)}")
        continue

    for entry in entries:
        message_id, values = entry[1][0]

        message_json = values[b'message']
        try:
            req_redis = json.loads(message_json)
        except json.JSONDecodeError as e:
            print(f"Error unmarshalling data: {str(e)}")
            continue

        print(f"Message: {req_redis}")
