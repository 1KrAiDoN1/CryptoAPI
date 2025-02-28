package handler

import (
	"fmt"
	"helloapp/internal/service"
	"helloapp/pkg/format"

	"helloapp/internal/database"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func HandleFunc() {
	http.HandleFunc("/home", showInfo) // Передаем данные о криптовалюте в обработчик
	http.HandleFunc("/crypto/", showCryptoDetails)
	http.HandleFunc("/registration", registration_window)
	http.HandleFunc("/sendUserRegistrationData", database.SendUserRegistrationData)
	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8080", nil)

}
func registration_window(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../pkg/templates/registration.html")
	if err != nil {
		log.Print("Ошибка при чтении шаблона:", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Print("Ошибка при выполнении шаблона:", err)
	}

}

func showInfo(w http.ResponseWriter, r *http.Request) {
	output, err := service.GetCryptoData()
	if err != nil {
		fmt.Println("Ошибка при получении данных", err, http.StatusInternalServerError)
		return
	}

	//fmt.Printf("Данные для шаблона: %+v\n", output)
	// Используем ParseFiles для загрузки шаблона из файла
	tmpl, err := template.New("home.html").Funcs(template.FuncMap{
		"formatLargeNumber": format.FormatLargeNumber, "formatLargeNumberForPercent": format.FormatLargeNumberForPercent, // Регистрируем функцию
	}).ParseFiles("../pkg/templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, output)
	if err != nil {
		http.Error(w, "Ошибка при обработке шаблона", http.StatusInternalServerError)
		return
	}
}

func showCryptoDetails(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID криптовалюты из URL
	id := strings.TrimPrefix(r.URL.Path, "/crypto/")
	if id == "" {
		http.Error(w, "ID криптовалюты не указан", http.StatusBadRequest)
		return
	}

	// Получаем данные о конкретной криптовалюте
	crypto, err := service.GetCryptoDataByID(id)
	if err != nil {
		http.Error(w, "Ошибка при получении данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Загружаем шаблон для страницы с деталями
	tmpl, err := template.New("crypto_details.html").Funcs(template.FuncMap{
		"formatLargeNumber": format.FormatLargeNumber, "formatLargeNumberForPercent": format.FormatLargeNumberForPercent, // Регистрируем функцию
	}).ParseFiles("../pkg/templates/crypto_details.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем шаблон с данными
	err = tmpl.Execute(w, crypto)
	if err != nil {
		http.Error(w, "Ошибка при обработке шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
