# Discord Bot - Weather and Dictionary Bot
This Discord bot provides weather information for cities and allows users to look up the meanings of English words. Additionally, it supports setting reminders with timers.

## Commands
!start: Start the bot and receive a greeting.

!help: View a list of available commands.

!weather [city]: Get the current weather in the specified city.

!define [word]: Look up the meaning of an English word.

!reminder [hh:mm:ss] [message]: Set a timer with a reminder message.

## Setup
1. Obtain a Discord bot token and replace YOUR_BOT_TOKEN with it in the main function.

2. Run the bot using go run main.go.

3. Invite the bot to your Discord server.


At first, i could not find an api that fetches data about the weather just by the name of the city, so in the code it is a little complicated as 
I programmed the bot to first find the latitude and longtitude of a city and only then fetch the weather based on the coordinates. Also, i could not find free API's that would translate text from one language to another 
so i just used a free api that gives a definition to some words in english





# Discord Бот - Погода и Словарь
Этот бот для Discord предоставляет информацию о погоде в городах и позволяет пользователям искать значения английских слов. Кроме того, он поддерживает установку напоминаний с таймерами.

## Команды
!start: Запустить бота и получить приветствие.

!help: Просмотреть список доступных команд.

!weather [город]: Получить текущую погоду в указанном городе.

!define [слово]: Найти значение английского слова.

!reminder [чч:мм:сс] [сообщение]: Установить таймер с сообщением-напоминанием.


## Настройка
1. Получите токен бота Discord и замените YOUR_BOT_TOKEN на него в функции main.

2. Запустите бота с помощью go run main.go.

3. Пригласите бота на ваш сервер Discord.



Изначально я не мог найти API, который предоставлял бы данные о погоде только по названию города, поэтому в коде немного усложнил процесс. Я настроил бота сначала на поиск широты и долготы города, а затем уже получение погодных данных на основе этих координат.
Также я не смог найти бесплатные API, которые предоставляли бы возможность перевода текста с одного языка на другой, поэтому я использовал бесплатный API, который предоставляет определения для некоторых слов на английском языке.



