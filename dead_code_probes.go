package tgbotapi

import "sync"

var dcProbe1 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-1") })
var dcProbe2 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-2") })
var dcProbe3 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-3") })
var dcProbe4 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-4") })
var dcProbe5 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-5") })
var dcProbe6 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-6") })
var dcProbe7 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-7") })
var dcProbe8 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-8") })
var dcProbe9 = sync.OnceFunc(func() { log.Printf("dead code probe hit probe=%s", "dc-telegram-bot-api-9") })
