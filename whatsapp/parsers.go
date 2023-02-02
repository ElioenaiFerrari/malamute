package whatsapp

import (
	"fmt"
	"strings"
)

func WhatsappPhone(phone string) string {
	return fmt.Sprintf("whatsapp:%s", phone)
}

func RawPhone(whatsappPhone string) string {
	return strings.ReplaceAll(whatsappPhone, "whatsapp:", "")
}
