package controllers

import (
	"apotek-management/config"
	"apotek-management/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetAllStok(c *gin.Context) {
	var stok []models.Stok
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	sortBy := c.DefaultQuery("sort_by", "nama_obat") 
	order := c.DefaultQuery("order", "asc")        

	log.Printf("Endpoint %s dipanggil oleh %s", c.FullPath(), c.ClientIP())

	if err := config.DB.Order(sortBy + " " + order).Limit(limit).Offset(offset).Find(&stok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data stok", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stok)
}

func CreateStok(c *gin.Context) {
	var input models.Stok
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format", "details": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

func UpdateStok(c *gin.Context) {
	var stok models.Stok
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := config.DB.First(&stok, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	var input models.Stok
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&stok).Updates(input)
	c.JSON(http.StatusOK, stok)
}

func DeleteStok(c *gin.Context) {
	var stok models.Stok
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := config.DB.First(&stok, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	if err := config.DB.Delete(&stok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func SearchStok(c *gin.Context) {
	var stok []models.Stok
	namaObat := c.DefaultQuery("nama_obat", "")

	if namaObat != "" {
		if err := config.DB.Where("nama_obat LIKE ?", "%"+namaObat+"%").Find(&stok).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stok)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama obat harus diisi"})
	}
}

func FilterStok(c *gin.Context) {
	var stok []models.Stok
	minHarga, minErr := strconv.ParseFloat(c.DefaultQuery("min_harga", "0"), 64)
	maxHarga, maxErr := strconv.ParseFloat(c.DefaultQuery("max_harga", "0"), 64)
	minJumlah, minErrJumlah := strconv.Atoi(c.DefaultQuery("min_jumlah", "0"))
	maxJumlah, maxErrJumlah := strconv.Atoi(c.DefaultQuery("max_jumlah", "0"))

	if minErr != nil || maxErr != nil || minErrJumlah != nil || maxErrJumlah != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter filter tidak valid"})
		return
	}

	query := config.DB
	if minHarga > 0 || maxHarga > 0 {
		query = query.Where("harga BETWEEN ? AND ?", minHarga, maxHarga)
	}

	if minJumlah > 0 || maxJumlah > 0 {
		query = query.Where("jumlah BETWEEN ? AND ?", minJumlah, maxJumlah)
	}

	if err := query.Find(&stok).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stok)
}

func GetStokSummary(c *gin.Context) {
	var totalJumlah int
	var avgHarga float64

	if err := config.DB.Model(&models.Stok{}).Select("SUM(jumlah)").Scan(&totalJumlah).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.Stok{}).Select("AVG(harga)").Scan(&avgHarga).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_stok":   totalJumlah,
		"average_harga": avgHarga,
	})
}
