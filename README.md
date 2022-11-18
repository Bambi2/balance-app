# HTTP API Микросервис для работы с балансом пользователей

## Запуск

```
make build && make run
```

## Примеры запросов

### Начисление денег на баланс пользователя или создание нового баланса с данной суммой, если пользователя раньше не было в системе 
<img width="1222" alt="image" src="https://user-images.githubusercontent.com/53175260/202730624-aec80852-4409-4357-ae44-0a4870d2327d.png">

### Получение баланса пользователя по его id
<img width="1227" alt="image" src="https://user-images.githubusercontent.com/53175260/202731285-a64f0b1e-719b-40cd-a960-2473661c710a.png">

### Резервирование средств с основного баланса на отдельном счете (ответ "ok", если все резервирование прошло успешно)
<img width="1224" alt="image" src="https://user-images.githubusercontent.com/53175260/202731883-1cf71f68-ece6-44f4-955f-5bc4737c2681.png">

### Если у пользователя не хватает денег на балансе, то приходит отказ в резервировании
<img width="1226" alt="image" src="https://user-images.githubusercontent.com/53175260/202732201-3302f13f-63f9-46cd-ac02-503f16e75a54.png">

### Метод признания выручки
<img width="1227" alt="image" src="https://user-images.githubusercontent.com/53175260/202732643-b385dee8-7563-4f5a-9792-5d7f01208578.png">

### Если на отдельном счете не осталось запрашиваемых денег для признания выручки, то приходит соответствующие сообщение с отказом
<img width="1225" alt="image" src="https://user-images.githubusercontent.com/53175260/202732861-84199fa8-54e0-4185-a2be-0a2ec93e4052.png">

### Метод разрезервирования денег, если остались не признанные для выручки деньги (после этого счет удаляется)
<img width="1222" alt="image" src="https://user-images.githubusercontent.com/53175260/202733531-7b3ccc26-e455-4d4a-89d1-45a07c01708a.png">

### Метод для получения месячного месячного отчета для бухгалтерии в CSV формате (на вход год и месяц), использовал сервис cloudinary
<img width="1223" alt="image" src="https://user-images.githubusercontent.com/53175260/202738226-2b045952-b3ac-495d-960a-2835a59c7f29.png">
<img width="1226" alt="image" src="https://user-images.githubusercontent.com/53175260/202738294-6ecf5d5a-ebe5-4a76-91d2-a6e59d13a06c.png">

### Метод для перечисления денег между пользователями
<img width="1226" alt="image" src="https://user-images.githubusercontent.com/53175260/202738896-f0bd09d0-78be-402a-b9d3-51a1cb627398.png">

### В случае, если сумма превышает баланс пользователя - отказ
<img width="1226" alt="image" src="https://user-images.githubusercontent.com/53175260/202739162-f5887406-70e6-4c28-9245-605a8ed79929.png">

### Метод получения списка комментариев транзакциий, отсортированных по дате и сумме (на вход id пользователя и limit с offset для пагинации)
<img width="1221" alt="image" src="https://user-images.githubusercontent.com/53175260/202740181-6436185a-ad3c-40b8-9c7b-7af51fe0e71f.png">

## Комментарий
- Вся деньги представлены в копейках, а не в рублях, чтобы не использовать float и не терять точность.
- Swagger документация лежит в директории docs и доступна по "http://localhost:8000/swagger/index.html"
- Я предположил, что услуги рекламные и признание выручки идет постепенно т.к. на вход идет сумма для признания. Поэтому пользователь может разрезервировать какие-то деньги, если они остались на отдельном счете и не были отправлены в выручку. 
