# Telegram Block and Card Bot

## Introduction

The **Telegram Block and Card Bot** is a simple and interactive Telegram bot designed to help users create and manage blocks of cards. Each block can contain multiple cards, making it easy to organize information for various purposes.

## Features

- **Create Blocks**: Users can create named blocks to categorize their cards.
- **Add Cards**: Users can easily add text-based cards to their created blocks.
- **View Blocks**: Display a list of all created blocks along with the count of cards in each.
- **View Cards**: Allows users to see all cards within a selected block.
- **Exit Block**: Option to exit from the card creation context back to the main menu.

## Getting Started

### Prerequisites

- Go programming language installed on your machine.
- A Telegram account.
- A bot token obtained from [BotFather](https://core.telegram.org/bots#botfather).

### Installation Steps

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/yourusername/telegram-block-card-bot.git
   cd telegram-block-card-bot

2. **Create Environment File:** Create a .env file in the project root directory and add your bot token:
   ```plaintxt
   BOT_TOKEN=your_bot_token_here

3. **Install Dependencies:** Run the following command to install the required Go packages:
   ```bash
   go get github.com/go-telegram-bot-api/telegram-bot-api/v5
   go get github.com/joho/godotenv

4. **Run the Bot:** Execute the following command to start the bot:
   ```bash
   go run main.go
   
### How to Use
1. Start the Bot: Open your Telegram app and find your bot using its username. Start the chat by sending the /start command.

2. Create a Block: Click on Создать блок (Create Block) and enter a name for your new block.

3. Add Cards: Once a block is created, select Добавить карточку (Add Card) and input the text for your card.

4. View Blocks: Select Показать блоки (Show Blocks) to view a list of all your blocks along with the number of cards in each.

5. View Cards in a Block: Choose Показать карточки блока (Show Cards in Block) and enter the block number to see all cards contained in that block.

6. Exit Block: Select Выйти из блока (Exit Block) to return to the main menu without losing your progress.

### What we can add?
1. Edit blocks 
2. Delete blocks or elements in card
3. Create database 