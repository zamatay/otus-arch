# otus-arch

### Для запуска проекта необходимо установить утилиту goose go install github.com/pressly/goose/v3/cmd/goose@latest

### Далее выполнить команду make run, что создаст БД и запустить все компоненты

### Для применения миграций выполнить команду run up

### ДЗ2

## SQL Скрипты
### Создание индекса
CREATE INDEX idx_first_last_name
ON public.users USING btree (first_name, last_name)
include (id, login, birthday, gender_id, city, enabled, interests);

### Анализирование запроса
explain analyse
select id, login, first_name, last_name, birthday, gender_id, city, enabled, interests
from users
where first_name like 'xxxxx%'  and last_name like 'xxxxx%'  
limit 100;


|                      |  Sample  |   Average |    Median |    90%   |     95%    |    99%  |  throughput  |
|----------------------|----------|-----------|-----------|----------|------------|---------|--------------|
| 1000 без индекса     |   2236   |  35514    |  45124    |  48452   |   48649    |   48787 |  20,6/sec    |
| 1000 с индексом      |   60360  |  929      |  844      |  2037    |   2339     |   4212  |  982,4/sec   |
| 100 без индекса      |   1287   |  4475     |  4932     |  5356    |   5423     |   5564  |  19,8/sec    |
| 100 с индексом       |   47562  |  76       |  18       |  208     |   238      |   309   |  1144/sec    |
| 10 без индекса       |   1209   |  460      |  415      |  797     |   854      |   974   |  20,1/sec    |
| 10 с индексом        |   56610  |  9        |  11       |  20      |   23       |   28    |  943/sec     |
