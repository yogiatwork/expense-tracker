## URLS

1. https://developers.facebook.com/tools/explorer?method=GET&path=1025639416765715%2Fphone_numbers&version=v26.0
2. https://business.facebook.com/latest/whatsapp_manager/phone_numbers?business_id=1025639416765715&asset_id=28308222772148929
3. https://developers.facebook.com/apps/2571852976579890/use_cases/customize/wa-configurations-v2/?use_case_enum=WHATSAPP_BUSINESS_MESSAGING&product_route=whatsapp-business&selected_tab=wa-configurations-v2&business_id=1025639416765715
4. https://medium.com/@hamzas2401/how-i-registered-my-whatsapp-business-number-on-meta-b175a290a451

## For Application Gin Development

1. https://www.deployhq.com/guides/gin
2. https://khimananda.com/blog/deploy-gin-to-production-a-practical-guide
3. https://tech-insider.org/gin-golang-tutorial-rest-api-2026/

## TODOs

1. Get the Secrets sorted and stored in a secure way, we only need user and private key from creds also think of using properties.
2. Spending User needs to be part of telegram message
3. Register the telegram bot command
4. Get User Expenses like /spend sapna should result in expense done from Sapna
5. Get Total Expenses /totalExpenses for that month
6. Get Balance /balance for the month which will be a balancesheet for total earned and total spend

## NGrok

Start the ngrok tunnel to expose the local server to the internet. This is necessary for Telegram to send messages to your local development server.

```bash
ngrok http 8080 --url https://landmass-second-doubling.ngrok-free.dev
```
