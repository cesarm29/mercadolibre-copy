package main

import (
	"fmt"
	"log"

	"marketplace/config"
	"marketplace/database"
	"marketplace/models"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()
	database.Connect(cfg)
	if database.DB == nil {
		log.Fatal("DB not connected")
	}

	var buyer models.User
	database.DB.Where("email = ?", "testbuyer@example.com").First(&buyer)
	if buyer.ID == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("Test1234"), bcrypt.DefaultCost)
		database.DB.Create(&models.User{
			Email:        "testbuyer@example.com",
			PasswordHash: string(hash),
			Name:         "Comprador Demo",
			Phone:        "3001112233",
			RoleID:       3,
			IsActive:     true,
		})
		fmt.Println("created buyer")
	}

	var seller models.User
	database.DB.Where("email = ?", "seller@example.com").First(&seller)
	if seller.ID == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("Seller1234"), bcrypt.DefaultCost)
		seller = models.User{
			Email:        "seller@example.com",
			PasswordHash: string(hash),
			Name:         "Vendedor Demo",
			Phone:        "3004445566",
			RoleID:       2,
			IsActive:     true,
		}
		database.DB.Create(&seller)
		fmt.Println("created seller id", seller.ID)
	}

	products := []models.Product{
		{Title: "iPhone 15 Pro 256GB Titanio", Description: "Telefono inteligente Apple, camara 48MP, chip A17 Pro. Nuevo, sellado.", Price: 5800000, Stock: 12, CategoryID: 1, Condition: "new", FreeShipping: true, Location: "Bogota", Images: []string{"https://picsum.photos/seed/iphone15/640/480"}},
		{Title: "Samsung Galaxy S24 Ultra 512GB", Description: "Camara 200MP, S Pen integrado, pantalla 6.8 pulgadas.", Price: 6200000, Stock: 8, CategoryID: 1, Condition: "new", FreeShipping: true, Location: "Medellin", Images: []string{"https://picsum.photos/seed/galaxys24/640/480"}},
		{Title: "Laptop ASUS ROG Strix G16 i7 16GB 512GB", Description: "Gaming, RTX 4060, pantalla 16 pulgadas 165Hz.", Price: 7800000, Stock: 5, CategoryID: 1, Condition: "new", FreeShipping: true, Location: "Cali", Images: []string{"https://picsum.photos/seed/rogstrix/640/480"}},
		{Title: "Nevera LG 425L No Frost Eficiente", Description: "Refrigerador inverter, dispensador de agua, 425 litros.", Price: 3900000, Stock: 7, CategoryID: 2, Condition: "new", FreeShipping: false, Location: "Bogota", Images: []string{"https://picsum.photos/seed/neveralg/640/480"}},
		{Title: "Lavadora Secadora Samsung 17kg", Description: "Carga frontal, 17kg, vaporizacion, conectividad SmartThings.", Price: 2900000, Stock: 9, CategoryID: 2, Condition: "new", FreeShipping: false, Location: "Bogota", Images: []string{"https://picsum.photos/seed/lavsam/640/480"}},
		{Title: "Horno Microondas Whirlpool 1.7 pies", Description: "1200W, grill, panel touch, 34L.", Price: 650000, Stock: 20, CategoryID: 2, Condition: "new", FreeShipping: false, Location: "Medellin", Images: []string{"https://picsum.photos/seed/microwh/640/480"}},
		{Title: "Tenis Nike Air Max 270 Hombre", Description: "Talla 42, amortiguacion Max Air, malla transpirable.", Price: 790000, Stock: 15, CategoryID: 3, Condition: "new", FreeShipping: true, Location: "Cali", Images: []string{"https://picsum.photos/seed/airmax270/640/480"}},
		{Title: "Chaqueta de Cuero Classic Biker", Description: "Cuero genuino, forro interior, tallas S-XL.", Price: 850000, Stock: 6, CategoryID: 3, Condition: "used", FreeShipping: true, Location: "Barranquilla", Images: []string{"https://picsum.photos/seed/jacketbiker/640/480"}},
		{Title: "Bicicleta Todo Terreno 26 Pulgadas", Description: "18 velocidades, frenos de disco, suspension delantera.", Price: 1200000, Stock: 10, CategoryID: 4, Condition: "new", FreeShipping: false, Location: "Bogota", Images: []string{"https://picsum.photos/seed/bicimtb/640/480"}},
		{Title: "Banda Elastica + Mancuernas 20kg Set", Description: "Kit completo de gimnasio en casa.", Price: 450000, Stock: 25, CategoryID: 4, Condition: "new", FreeShipping: true, Location: "Medellin", Images: []string{"https://picsum.photos/seed/gymkit/640/480"}},
		{Title: "Tabla de Picar + Set de Ollas", Description: "Set cocina 12 piezas, antiadherente.", Price: 380000, Stock: 18, CategoryID: 5, Condition: "new", FreeShipping: true, Location: "Cali", Images: []string{"https://picsum.photos/seed/setcocina/640/480"}},
		{Title: "Taladro Percutor Inalambrico 20V", Description: "Kit con 2 baterias, maletin, 20 accesorios.", Price: 720000, Stock: 14, CategoryID: 7, Condition: "new", FreeShipping: true, Location: "Bogota", Images: []string{"https://picsum.photos/seed/taladro/640/480"}},
	}

	added := 0
	for _, p := range products {
		p.SellerID = seller.ID
		var existing models.Product
		if err := database.DB.Where("title = ?", p.Title).First(&existing).Error; err == nil {
			fmt.Printf("existing: %s\n", p.Title)
			continue
		}
		if err := database.DB.Create(&p).Error; err != nil {
			fmt.Printf("CREATE ERROR [%s]: %v\n", p.Title, err)
			continue
		}
		fmt.Printf("created product: %s (id=%d)\n", p.Title, p.ID)
		added++
	}
	fmt.Printf("products checked/created: %d\n", added)

	var count int64
	database.DB.Model(&models.Product{}).Count(&count)
	fmt.Printf("total products in DB: %d\n", count)
}