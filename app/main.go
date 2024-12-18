package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Text string
}

type Block struct {
	Name  string
	Cards []Card
}

var blocks = make(map[int64][]Block)

func GetToken() string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	return os.Getenv("BOT_TOKEN")
}

func handleStart(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Выберите действие:")
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Создать блок")},
		},
		ResizeKeyboard: true,
	}
	bot.Send(msg)
}

func handleCreateBlock(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Введите название блока:")
	bot.Send(msg)
}

func handleAddCard(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Введите текст карточки:")
	bot.Send(msg)
}

func showBlocks(bot *tgbotapi.BotAPI, chatID int64) {
	userBlocks := blocks[chatID]
	if len(userBlocks) == 0 {
		msg := tgbotapi.NewMessage(chatID, "У вас пока нет созданных блоков.")
		bot.Send(msg)
		return
	}

	var blocksList []string
	for i, block := range userBlocks {
		blocksList = append(blocksList, fmt.Sprintf("%d. %s (%d карточек)", i+1, block.Name, len(block.Cards)))
	}

	msg := tgbotapi.NewMessage(chatID, "Ваши блоки:\n"+strings.Join(blocksList, "\n"))
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{tgbotapi.NewKeyboardButton("Создать блок")},
			{tgbotapi.NewKeyboardButton("Показать карточки блока")},
		},
		ResizeKeyboard: true,
	}
	bot.Send(msg)
}

func showCardsInBlock(bot *tgbotapi.BotAPI, chatID int64, blockIndex int) {
	userBlocks := blocks[chatID]
	if blockIndex < 0 || blockIndex >= len(userBlocks) {
		msg := tgbotapi.NewMessage(chatID, "Блок не найден.")
		bot.Send(msg)
		return
	}

	block := userBlocks[blockIndex]
	if len(block.Cards) == 0 {
		msg := tgbotapi.NewMessage(chatID, "В этом блоке пока нет карточек.")
		bot.Send(msg)
		return
	}

	var cardsList []string
	cardsList = append(cardsList, fmt.Sprintf("Карточки блока \"%s\":", block.Name))
	for i, card := range block.Cards {
		cardsList = append(cardsList, fmt.Sprintf("%d. %s", i+1, card.Text))
	}

	msg := tgbotapi.NewMessage(chatID, strings.Join(cardsList, "\n"))
	bot.Send(msg)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(GetToken())
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	var creatingBlock bool
	var currentBlock Block
	var addingCard bool
	var waitingForBlockNumber bool

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messageText := update.Message.Text

		switch {
		case messageText == "/start":
			handleStart(bot, chatID)
			creatingBlock = false
			addingCard = false
			waitingForBlockNumber = false

		case messageText == "Создать блок":
			handleCreateBlock(bot, chatID)
			creatingBlock = true
			addingCard = false
			waitingForBlockNumber = false

		case messageText == "Показать блоки":
			showBlocks(bot, chatID)
			creatingBlock = false
			addingCard = false
			waitingForBlockNumber = false

		case messageText == "Показать карточки блока":
			msg := tgbotapi.NewMessage(chatID, "Введите номер блока:")
			bot.Send(msg)
			waitingForBlockNumber = true
			creatingBlock = false
			addingCard = false

		case waitingForBlockNumber:
			blockNum, err := strconv.Atoi(messageText)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите корректный номер блока.")
				bot.Send(msg)
				continue
			}
			showCardsInBlock(bot, chatID, blockNum-1)
			waitingForBlockNumber = false

		case creatingBlock:
			currentBlock = Block{
				Name:  messageText,
				Cards: []Card{},
			}
			blocks[chatID] = append(blocks[chatID], currentBlock)
			creatingBlock = false
			addingCard = true

			msg := tgbotapi.NewMessage(chatID, "Блок создан! Теперь вы можете добавлять карточки или выйти из блока.")
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
				Keyboard: [][]tgbotapi.KeyboardButton{
					{tgbotapi.NewKeyboardButton("Добавить карточку")},
					{tgbotapi.NewKeyboardButton("Выйти из блока")},
				},
				ResizeKeyboard: true,
			}
			bot.Send(msg)

		case addingCard && messageText == "Добавить карточку":
			handleAddCard(bot, chatID)

		case addingCard && messageText != "Выйти из блока" && messageText != "Показать блоки":
			if len(blocks[chatID]) > 0 {
				lastBlockIndex := len(blocks[chatID]) - 1
				card := Card{Text: messageText}
				blocks[chatID][lastBlockIndex].Cards = append(blocks[chatID][lastBlockIndex].Cards, card)
				msg := tgbotapi.NewMessage(chatID, "Карточка добавлена! Можете добавить ещё или выйти из блока.")
				bot.Send(msg)
			}

		case messageText == "Выйти из блока":
			msg := tgbotapi.NewMessage(chatID, "Вы вышли из блока. Выберите действие:")
			msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
				Keyboard: [][]tgbotapi.KeyboardButton{
					{tgbotapi.NewKeyboardButton("Создать блок")},
					{tgbotapi.NewKeyboardButton("Показать блоки")},
				},
				ResizeKeyboard: true,
			}
			bot.Send(msg)
			creatingBlock = false
			addingCard = false
			waitingForBlockNumber = false

		default:
			if !waitingForBlockNumber && !creatingBlock && !addingCard {
				msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
				msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
					Keyboard: [][]tgbotapi.KeyboardButton{
						{tgbotapi.NewKeyboardButton("Создать блок")},
						{tgbotapi.NewKeyboardButton("Показать блоки")},
					},
					ResizeKeyboard: true,
				}
				bot.Send(msg)
			}
		}
	}
}
