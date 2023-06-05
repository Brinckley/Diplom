package telegram

const msgUnknownCommand = "Неизвестная команда 😞🥁!" //"Unknown command 😞🥁!"
const msgStartCommand = "Привет 👨‍💻👨‍💻👨‍💻\n" +
										//"Hello! 👨‍💻👨‍💻👨‍💻\n" +
										msgListCommand
const msgSuccessfulRegistration = "Вы успешно зарегестрированы! \U0001FAF6" //"You've registered successfully! \U0001FAF6"
const msgHelpCommand = "🧙‍Список команд.... \n" + msgListCommand

// "🧙‍Some useful instructions right here.... \n\n️" + msgListCommand
const msgUserExistsCommand = "😾😾😾Пользователя с таким именем не существует!😾😾😾" //"😾😾😾User with such name already exists!😾😾😾"
const msgAddArtist = ""

const msgSearchArtist = "Введите имя артиста 👨‍🎤 для поиска: "

const msgNoArtists = "🚨🚨🚨Артист не найден!!!🚨🚨🚨"         //"🚨🚨🚨No artists in the found!!!🚨🚨🚨"
const msgNoFavorites = "😱😱😱Нет артистов в избранном!😱😱😱" //"😱😱😱You have no favorite artists!😱😱😱"
const msgListCommand =
// "Get Help by pressing /help 🧙‍♂️\n" +
"Помощь -> /help 🧙‍♂️\n" +
	"Меню -> /catalog 📄\n" +
	"Подписаться -> /subscribe 🖋\n" +
	"Отписаться -> /unsubscribe ❌\n" +
	"Избранное -> -> /favorites 🎧\n"

//"[DEBUG] If you are not registered in the service, please do it by pressing -> /register 🖖\n" +
//"View the menu -> /catalog 📄\n" +
//"View catalog by pressing -> /catalog 🎸\n" +
//"Add the artist to favorites by pressing -> /subscribe 🖋\n" +
//"Remove the artist from favorites by pressing -> /unsubscribe ❌\n" +
//"Get your favorites here -> /favorites 🎧\n"

const msgAskArtist = "Haven't found artist you want? Add it here -> /addArtist"

const msgIntroArtists = "Вот список имеющихся артистов: \n\n" //"Here is the complete list of artists:\n\n"
const msgIntroFavorites = "Избранные артисты:\n\n"            //"Here is the list of your favorite artists:\n\n"

const msgSubscribeQuestion = "Введите имя артиста 👨‍🎤 для добавления в избранное🔥..."    //"Please, enter the name of the 👨‍🎤 artist you want to 🖋add🖋 to favorites🔥..."
const msgSubscribeSuccess = "Вы успешно подписались на артиста! 🎷"                       //"You've successfully 🎷🎷🎷 subscribed to the news of artist"
const msgSubscribeFail = "😢😢😢 Артист с таким именем не найден в базе 😢😢😢"                //"😢😢😢 Artist with such name not found 😢😢😢"
const msgSubscribeAlready = "🕶Уже в избранном🕶"                                          //"🕶Already subscribed!🕶"
const msgUnsubscribeQuestion = "Введите имя артиста 👨‍🎤 для исключения ❌ из избранного🔥" //"Please, enter the name of the 👨‍🎤 artist to ❌remove❌ from favorites🔥..."
const msgUnsubscribeFail = "Что-то пошло не так 👺. Вы уверены, что артист в избранном🔥?" //"Something went wrong👺👺👺. Is there such artist in your favorites🔥?"
const msgUnsubscribeSuccess = "Артист удален из избранного!🧟"                            //"You've successfully 🧟 unsubscribed to the news of artist"

const msgStartAgainCommand = "С возвращением!" //"Welcome back 🥃"

const msgRandomMessage = "Неправильное сообщение 🤖"                     //"Invalid message 🤖"
const msgRandomMessageFinal = "Прекратите писать неверные сообщения ☠️" //"Stop writing useless messages! ☠️"
