package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	uniqueVisitors = make(map[int64]struct{}) // Для хранения уникальных ID
	mu             sync.Mutex                 // Мьютекс для безопасного доступа
)

func resetUniqueVisitors() {
	for {
		time.Sleep(24 * time.Hour) // Ждём сутки
		mu.Lock()
		uniqueVisitors = make(map[int64]struct{}) // Очищаем карту
		log.Println("Ежедневная статистика уникальных посещений сброшена.")
		mu.Unlock()
	}
}
func handleStart(userID int64) {
	mu.Lock()
	defer mu.Unlock()

	// Проверяем, есть ли пользователь в карте
	if _, exists := uniqueVisitors[userID]; !exists {
		uniqueVisitors[userID] = struct{}{} // Добавляем пользователя
		log.Printf("Новый уникальный посетитель: %d. Всего уникальных: %d\n", userID, len(uniqueVisitors))
	}
}

func main() {
	go resetUniqueVisitors()
	_ = godotenv.Load(".env")

	// Получение токена из переменной окружения
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	// Инициализация бота
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	bot.Debug = true // Включить режим отладки

	// Настраиваем обновления от бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	// Открываем файл для записи логов (или создаём, если его нет)
	logFile, err := os.OpenFile("visitors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл для логов: %v", err)
	}
	// Настраиваем логгер на запись в файл
	log.SetOutput(logFile)

	// Закрываем файл при завершении программы
	defer logFile.Close()

	// Остальной код
	go resetUniqueVisitors()
	for update := range updates {
		if update.Message != nil { // Если пришло сообщение
			switch update.Message.Text {
			case "/start":
				handleStart(update.Message.From.ID)
				// Создаём кнопки с подходящими смайликами и описаниями
				buttons := [][]tgbotapi.InlineKeyboardButton{
					{
						tgbotapi.NewInlineKeyboardButtonData("🧘‍♀️ Это всё я", "practice_1"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("🦸‍♂️ Сила больше боли", "practice_2"),
						tgbotapi.NewInlineKeyboardButtonData("🌏 Я пришел жить вслух", "practice_3"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("😟 Разговор с тревогой", "practice_4"),
						tgbotapi.NewInlineKeyboardButtonData("🌌 Мир уже состоит из твоих смыслов", "practice_5"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("🌱 Начни жить свою жизнь", "practice_6"),
						tgbotapi.NewInlineKeyboardButtonData("🧘‍♂️ Сравнивая себя, ты теряешь себя", "practice_7"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("🌼 Рай под ногами", "practice_8"),
					},
				}

				// Создаём клавиатуру с кнопками
				keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

				// Формируем сообщение с кнопками
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Добро пожаловать в проект "Зеркала"! Наши психологи здесь делятся практиками, которые помогают в жизни им самим. Почувствуй себя как дома и выбери тему для проработки:`)
				msg.ReplyMarkup = keyboard

				bot.Send(msg)

			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не распознана. Введите /start для начала.")
				bot.Send(msg)
			}
		} else if update.CallbackQuery != nil { // Если пользователь нажал на кнопку
			var responseText string
			switch update.CallbackQuery.Data {
			case "practice_1":
				responseText = "*Когда я замечаю, что теряю себя, смущаюсь или боюсь чьей-то реакции и из-за этого становлюсь наигранной и напряженной, я напоминаю себе: я есть и этого достаточно.*\n\n" +
					"Мне помогает дыхательная техника Треугольник силы.\n\n" +
					"Я кладу одну ладонь на живот, а вторую на крестец и чувствую тепло между ними, в животе.\n" +
					"С выдохом я мысленно направляю это тепло в правую ногу и пускаю эту энергию по правой ноге до самой земли и дальше, через ступню, под землю.\n" +
					"Под землей я направляю эту энергию в левую стопу и со вдохом представляю, как поднимаю ее обратно, в живот.\n" +
					"Так, я дышу несколько кругов, заземляюсь и чувствую себя до тех пор, пока не пойму: я есть и этого достаточно.\n\n" +
					"— _Эту практику создала Марта Куклина. Психолог, гештальттерапевт._\n" +
					"[Сайт для связи](https://kuklina.pro)"

				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown" // Указываем режим Markdown для форматирования текста
				bot.Send(msg)

			case "practice_2":
				responseText = "*Когда у меня опускаются руки или что-то выбивает почву из-под ног – я напоминаю себе, что уже многое смогла и многого достигла. Я вспоминаю о своей внутренней силе. Чтобы не забывать о ней – я иногда возвращаюсь к этому упражнению.*\n\n" +
					"**Поиск и присвоение сил.**\n\n" +
					"*Шаг 1.*  \n                    Я предлагаю тебе написать список ресурсов, начиная с самого раннего детства в контексте того, что помогло выжить и помогает жить. Для этого надо повспоминать сложные периоды/ситуации из жизни, от последних, настоящих до самого детства, и выписать ответ на вопрос:  \n                    _Что мне тогда помогло, что мной двигало, благодаря каким своим способностям я выжила и справилась?_  \n                    Это могут быть и одобряемые вещи, и социально неодобряемые (обманул, изменила и пр.). В работе делай акцент внутри себя на том, что какой бы ни был способ — \"главное, что я справилась\".  \n                    В этом шаге выписывай без редактуры, даже если внутренний критик будет говорить \"тоже мне таланты!\". Просто пиши всё, что помогло. Постарайся не меньше 20 пунктов.\n\n" +
					"*Шаг 2.*  \n                    Подели список на внешние и внутренние ресурсы.  \n                    _Внутренние:_ это «я сам/а» умею, сделал, поняла.  \n                    _Внешние:_ \"повезло, помогли, совпало, хорошие друзья, добрые люди, удачный был день\".  \n                    Внутренние – пока просто отложи в сторону. Внешние конвертируй во внутренние.  \n                    Например: _мне помогли люди_ = _я умею попросить помощи у правильных людей_ ИЛИ _я умею поддерживать отношения_.\n\n" +
					"*Шаг 3. Конвертация.*  \n                    В этом шаге соедини внутренние и внешние, которые конвертировала во внутренние – это список твоих внутренних сил. В нём поменяй те формулировки, которые тебе не нравятся, так, чтобы они тебе нравились, чтобы ты с гордостью говорила это про себя. Меняй до тех пор, пока не будет приятно говорить так о себе.\n\n" +
					"*Шаг 4.*  \n                    У тебя есть список (обязательно его написать или напечатать в заметки/документ, не \"в уме\") из твоих способностей, особенностей, навыков, талантов, которые тебе помогают в сложных ситуациях. Все внутренние (начинаются с \"я\"/\"мне\"). Желательно, чтобы при его прочтении \"выравнивалась спинка\".\n\n" +
					"*Шаг 5. Телесное закрепление.*  \n                    Надиктуй этот список себе в аудио формате: прочитай все пункты в формулировке \"я умею...\" с небольшой паузой после каждого пункта. Слушай эту запись во время упадка сил, обращай внимание на телесные реакции во время прослушивания.\n\n" +
					"_ЗАДАЧА:_  \n                    Закрепить телесное состояние от знания своих сил в самых маленьких ощущениях (\"накачать попу\"). Чтобы когда оно наработается – уметь вызывать это состояние не мыслями, а ощущениями, мышцами, дыханием. То есть, ты как бы выучиваешь \"асану\" и потом будешь не думать о себе хорошее, а просто телом \"вставать\" в твою асану :)\n\n" +
					"**Сила внутри тебя.**\n\n" +
					"— _Эту практику создала Марта Куклина. Психолог, гештальттерапевт._  \n                    Сайт для связи: [https://kuklina.pro](https://kuklina.pro)"

				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "practice_3":
				responseText = "*Одно из самых больших разочарований в жизни - это прожить не свою жизнь. А что останавливает человека выбирать свой собственный, настоящий, тот самый путь?*\n\n" +
					"Ниже предлагаю несколько фраз, которые могут помочь выйти из привычных паттернов.\n\n" +
					"Попробуйте начать рассуждения на тему **“Если бы я был посмелее, я бы тогда..”** и напишите минимум 10 пунктов. Если захочется больше, не ограничивайте себя.\n\n" +
					"Далее новый список **“Если бы я себя любил, я бы…”** и тоже минимум 10 пунктов.\n\n" +
					"Ну и последнее **“Если бы мне было можно жить так, как мне хочется, я бы тогда…”**\n\n" +
					"*Эту практику создала Федорец Екатерина. Психолог, гештальт и логотерапевт.*  \n" +
					"Сайт для связи: [https://fedorets-psy.ru/](https://fedorets-psy.ru/)"

				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "practice_4":
				responseText = "Описание практики 'Разговор с тревогой'..."

				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "practice_5":
				responseText = "„Нет такой ситуации, в которой нам не была бы предоставлена жизнью возможность найти смысл, и нет такого человека, для которого жизнь не держала бы наготове какое-нибудь дело. Возможность осуществить смысл всегда уникальна, и человек, который может ее реализовать, всегда неповторим“  \n        — Виктор Франкл - австрийский психиатр, психолог, философ и невролог, бывший узник нацистского концентрационного лагеря.  \n        Автор мирового бестселлера «Сказать жизни „Да!“: Психолог в концлагере».\n\nЕсли представить, что внутри тебя спрятан компас, который всегда показывает тебе направление твоего предназначения и твоего смысла. Куда бы этот компас показал сейчас? В какой области был твой “север”? Куда в течение всей жизни тебя тянет?\n\nЭту практику создала Федорец Екатерина. Психолог, гештальт и логотерапевт.  \nСайт для связи: [https://fedorets-psy.ru/](https://fedorets-psy.ru/)"
				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "practice_6":
				responseText = "О чем чаще всего жалеют люди перед смертью?  \n" +
					"Бронни Вэр называют самой известной паллиативной медсестрой в мире. Она долгие годы работала в хосписе и провела множество бесед с людьми, которым оставалось всего несколько дней. Она заметила, что пациенты сожалеют о похожих вещах.  \n" +
					"И чаще всего люди жалеют о том, что им не хватило смелости жить так, как им хочется.\n\n" +
					"Все мы в той или иной степени зависим от мнения окружающих и хотим их одобрения. Но к сожалению, это также приводит к тому, что люди отодвигают/замалчивают свои истинные желания, чтобы быть хорошими/удобными для других. И в итоге проживают не свою жизнь.  \n" +
					"Страшно оказаться без поддержки, страшно, что тебя осудят, от тебя отвернуться, если ты начнёшь поступать так, как тебе хочется.\n\n" +
					"Я тоже с этим сталкиваюсь и мне тоже бывает страшно что-то делать. Появляется ощущение, что сейчас все резко обратят на меня внимание и начнут осуждать.  \n" +
					"Что делать?  \n\n" +
					"Важно помнить, что все люди заняты собой, своей жизнью, они большую часть времени думают как и вы о себе. Им совершенно нет дела до вас.  \nИ даже если кто-то вас осуждает, вы не обязаны вести себя так, как им нравится.  \n" +
					"В моменты, когда вы хотите пожертвовать своими желаниями ради того, чтобы угодить другим, задайте себе эти вопросы:  \n\nЧто я теряю, когда выбираю вместо себя желания другого человека?  \nК чему это обычно приводит?  \nКак долго я еще готов/а это делать?\n\n" +
					"Контакты:  \nЕкатерина Жемлаускас  \n[Telegram](https://t.me/zhempsy)  \n[Instagram](https://www.instagram.com/ekaterina_zhemlauskas)"
				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "practice_7":
				responseText = "«Когда я смотрю на людей в своей сфере и начинаю себя сравнивать с ними - я впадаю в уныние. Они слишком далеко от меня, у них все легко и естественно получается, наверное, со мной что-то не так. Я какая-то не такая». Знакомо?  \n\n" +
					"Абсолютно каждый из нас хотя бы раз сравнивал себя с другими. Когда мы себя с кем-то сравниваем, мы теряем себя, мы теряем свою уникальность.\n\n" +
					"Ты открываешь соц. сети, и видишь, как твоя знакомая Маша отдыхает в Дубае и зарабатывает миллионы. В этот момент ты, скорее всего, думаешь: «Я никчемная, я ничего не могу, вот Машка да, Машка молодец, а чего добилась я? Ничего, и впадаешь в страдания».\n\n" +
					"Мне знакомо это состояние. Я тоже бываю в таком состоянии. Чаще всего оно приходит тогда, когда у меня что-то не получается. Мне кажется, что я неудачница и все, что я делаю – бесполезно.\n\n" +
					"Выйти из этого состояния помогает следующее упражнение.\n\n" +
					"1. Я сажусь и подробно отвечаю на следующие вопросы:  \n- Какую цену я плачу, когда сравниваю себя с другими?  \n- Что я теряю, когда нахожусь в этом состоянии?  \n\n2. Я смотрю на ответы, которые получились. Дальше я принимаю решение, что я хочу с этим сделать? Остаться в этом состоянии «сравнения и уныния» или принять решение выйти из него и пойти жить свою жизнь.\n\n" +
					"Автор: Екатерина Жемлаускас  \nСайт для связи: [https://t.me/zhempsy](https://t.me/zhempsy)  \n[Instagram](https://www.instagram.com/ekaterina_zhemlauskas)"
				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "practice_8":
				responseText = "Мой мир не сахар по жизни, особенно раньше. Проваливался в сверх тревожные и навязчивые состояния, в которых хочется уйти, сбежать.  \nОт ума, который объяснял, что мир не для меня, до эмоций и телесных ощущений, которые гнали домой.  \n\n" +
					">Мне помогает техника обнуления этих ощущений.  \n\n" +
					"Спрашиваю себя: Кто это думает?  \nОтвечаю: Я! А кто я?  \nПосле этого ловлю состояние момента, дыхания  \nИ повторяю это, когда приходят мысли и эмоции.  \nВыходя на здесь и сейчас, а там только принятие, где -  \nРай под ногами моими.\n\n" +
					"Автор: Алексей Сахаров"
				// Создание кнопки "Назад"
				backButton := tgbotapi.NewInlineKeyboardButtonData("Назад", "menu")
				backKeyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(backButton))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, responseText)
				msg.ReplyMarkup = backKeyboard
				msg.ParseMode = "Markdown"
				bot.Send(msg)

			case "menu":
				// Создаём кнопки с подходящими смайликами и описаниями
				buttons := [][]tgbotapi.InlineKeyboardButton{
					{
						tgbotapi.NewInlineKeyboardButtonData("🧘‍♀️ Это всё я", "practice_1"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("🦸‍♂️ Сила больше боли", "practice_2"),
						tgbotapi.NewInlineKeyboardButtonData("🌏 Я пришел жить вслух", "practice_3"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("😟 Разговор с тревогой", "practice_4"),
						tgbotapi.NewInlineKeyboardButtonData("🌌 Мир уже состоит из твоих смыслов", "practice_5"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("🌱 Начни жить свою жизнь", "practice_6"),
						tgbotapi.NewInlineKeyboardButtonData("🧘‍♂️ Сравнивая себя, ты теряешь себя", "practice_7"),
					},
					{
						tgbotapi.NewInlineKeyboardButtonData("🌼 Рай под ногами", "practice_8"),
					},
				}

				// Создаём клавиатуру с кнопками
				keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

				// Формируем сообщение с кнопками
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, `Добро пожаловать в проект "Зеркала"! Наши психологи здесь делятся практиками, которые помогают в жизни им самим. Почувствуй себя как дома и выбери тему для проработки:`)
				msg.ReplyMarkup = keyboard

				bot.Send(msg)
			}
		}
	}
}
