package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Produit struct {
	Id               int
	Nom              string
	Description      string
	Prix             float64
	Reduction        float64
	Image            string
	Lareduc          bool
	PrixReduit       float64
	PourcentageReduc int
}

var produits = []Produit{
	{
		Id:               1,
		Nom:              "PALACE PULL A CAPUCHE UNISEXE CHASSEUR",
		Description:      "Pull unisexe confortable",
		Prix:             129.99,
		Reduction:        0.20,
		Image:            "/static/products/19A.webp",
		Lareduc:          true,
		PrixReduit:       103.99,
		PourcentageReduc: 20,
	},
	{
		Id:               2,
		Nom:              "PALACE PULL A CAPUCHON MARINE",
		Description:      "Pull marine stylé",
		Prix:             119.00,
		Reduction:        0.10,
		Image:            "/static/products/21A.webp",
		Lareduc:          true,
		PrixReduit:       107.10,
		PourcentageReduc: 10,
	},
	{
		Id:          3,
		Nom:         "PALACE PULL CREW PASSEPOSE NOIR",
		Description: "Pull noir classique",
		Prix:        99.50,
		Image:       "/static/products/22A.webp",
		Lareduc:     false,
	},
	{
		Id:               4,
		Nom:              "PALACE WASHED TERRY 1/4 PLACKET HOOD MOJITO",
		Description:      "Hoodie vert mojito",
		Prix:             139.00,
		Reduction:        0.15,
		Image:            "/static/products/16A.webp",
		Lareduc:          true,
		PrixReduit:       118.15,
		PourcentageReduc: 15,
	},
	{
		Id:               5,
		Nom:              "PALACE PANTALON BOSSY JEAN STONE",
		Description:      "Jean stone coupe bossy",
		Prix:             149.90,
		Reduction:        0.05,
		Image:            "/static/products/34B.webp",
		Lareduc:          true,
		PrixReduit:       142.41,
		PourcentageReduc: 5,
	},
	{
		Id:               6,
		Nom:              "PALACE PANTALON CARGO GORE-TEX R-TEK NOIR",
		Description:      "Cargo Gore-Tex noir",
		Prix:             199.00,
		Reduction:        0.25,
		Image:            "/static/products/33B.webp",
		Lareduc:          true,
		PrixReduit:       149.25,
		PourcentageReduc: 25,
	},
}

func main() {
	os.MkdirAll("./assets/products", os.ModePerm)
	os.MkdirAll("./src/templates", os.ModePerm)

	temp, err := template.ParseGlob("./src/templates/*.html")
	if err != nil {
		fmt.Println("Erreur template:", err)
		os.Exit(1)
	}

	http.Handle("/src/templates/", http.StripPrefix("/src/templates/", http.FileServer(http.Dir("./src/templates"))))

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		err := temp.ExecuteTemplate(w, "home", produits)
		if err != nil {
			http.Error(w, "Erreur Templates", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/produit", func(w http.ResponseWriter, r *http.Request) {
		idParam := r.URL.Query().Get("id")
		produitId, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		for _, product := range produits {
			if product.Id == produitId {
				err := temp.ExecuteTemplate(w, "produit", product)
				if err != nil {
					http.Error(w, "Erreur Templates", http.StatusInternalServerError)
				}
				return
			}
		}

		http.Error(w, "Produit non trouvé", http.StatusNotFound)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ajouter un produit - StreetShop</title>
    <link rel="stylesheet" href="/src/templates/style.css">
</head>
<body>
    <header class="nav-link">
        <img class="logo" src="/assets/img/logo/1.png" alt="StreetShop Logo">
        <nav>
            <a href="/home">Accueil</a>
            <a href="/add">Add</a>
        </nav>
    </header>

    <div class="container">
        <h2 class="form-title">Ajouter un nouveau produit</h2>
        <form action="/process-add" method="POST">
            <div class="form-group">
                <label for="nom">Nom du produit:</label>
                <input type="text" id="nom" name="nom" required>
            </div>
            
            <div class="form-group">
                <label for="description">Description:</label>
                <textarea id="description" name="description" required></textarea>
            </div>
            
            <div class="form-group">
                <label for="prix">Prix:</label>
                <input type="number" id="prix" name="prix" step="0.01" min="0" required>
            </div>
            
            <div class="form-group">
                <label for="reduction">Réduction (0-1, 0 = pas de réduction):</label>
                <input type="number" id="reduction" name="reduction" step="0.01" min="0" max="1" value="0">
            </div>
            
            <div class="form-group">
                <label for="image">URL de l'image:</label>
                <input type="text" id="image" name="image" value=/assets/img/products/" required>
            </div>
            
            <button type="submit" class="submit-button">Ajouter le produit</button>
        </form>
        
        <a href="/home" class="back-button">Retour à la liste</a>
    </div>
</body>
</html>
`)
	})
	http.HandleFunc("/process-add", func(w http.ResponseWriter, r *http.Request) {
		nom := r.FormValue("nom")
		description := r.FormValue("description")
		prixStr := r.FormValue("prix")
		reductionStr := r.FormValue("reduction")
		image := r.FormValue("image")

		if nom == "" || description == "" || prixStr == "" {
			http.Error(w, "Tous les champs requis doivent être remplis", http.StatusBadRequest)
			return
		}

		prix, err := strconv.ParseFloat(prixStr, 64)
		if err != nil || prix < 0 {
			http.Error(w, "Prix invalide", http.StatusBadRequest)
			return
		}

		reduction, err := strconv.ParseFloat(reductionStr, 64)
		if err != nil || reduction < 0 || reduction > 1 {
			http.Error(w, "Réduction invalide (doit être entre 0 et 1)", http.StatusBadRequest)
			return
		}

		prixReduit := prix
		if reduction > 0 {
			prixReduit = prix * (1 - reduction)
		}

		newId := 1
		if len(produits) > 0 {
			newId = produits[len(produits)-1].Id + 1
		}

		nouveauProduit := Produit{
			Id:               newId,
			Nom:              nom,
			Description:      description,
			Prix:             prix,
			Reduction:        reduction,
			Image:            image,
			Lareduc:          reduction > 0,
			PrixReduit:       prixReduit,
			PourcentageReduc: int(reduction * 100),
		}

		produits = append(produits, nouveauProduit)

		http.Redirect(w, r, fmt.Sprintf("/produit?id=%d", newId), http.StatusSeeOther)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
