# Добавление имени БД
curl --request PUT --data example_app localhost:8500/v1/kv/example.app/db/name

# Добавление имени пользователя
curl --request PUT --data postgres localhost:8500/v1/kv/example.app/db/username

# Добавление пароля
curl --request PUT --data example_pass localhost:8500/v1/kv/example.app/db/password