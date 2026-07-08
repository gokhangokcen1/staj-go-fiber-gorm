<script setup>
import { ref, onMounted } from 'vue'
import OgrenciService from './services/OgrenciService'
import OgrenciForm from './components/OgrenciForm.vue'
import OgrenciList from './components/OgrenciList.vue'

const ogrenciler = ref([])
const duzenlenecekOgrenci = ref(null)

async function ogrencileriGetir() {
  try {
    const response = await OgrenciService.getAll()
    ogrenciler.value = response.data
  } catch (error) {
    console.error('Ogrenciler getirilemedi:', error)
  }
}

async function ogrenciKaydet(ogrenciVerisi) {
  try {
    if (duzenlenecekOgrenci.value) {
      await OgrenciService.update(duzenlenecekOgrenci.value.ID, ogrenciVerisi)
      duzenlenecekOgrenci.value = null
    } else {
      await OgrenciService.create(ogrenciVerisi)
    }
    await ogrencileriGetir()
  } catch (error) {
    console.error('Kaydetme basarisiz:', error)
  }
}

async function ogrenciSil(id) {
  try {
    await OgrenciService.delete(id)
    await ogrencileriGetir()
  } catch (error) {
    console.error('Silme basarisiz:', error)
  }
}

function duzenlemeyeBasla(ogrenci) {
  duzenlenecekOgrenci.value = ogrenci
}

onMounted(() => {
  ogrencileriGetir()
})
</script>

<template>
  <div class="container">
    <h1>Öğrenci Yönetim Sistemi</h1>

    <OgrenciForm
      :duzenlenecek="duzenlenecekOgrenci"
      @kaydet="ogrenciKaydet"
    />

    <OgrenciList
      :ogrenciler="ogrenciler"
      @sil="ogrenciSil"
      @duzenle="duzenlemeyeBasla"
    />
  </div>
</template>

<style>
.container {
  max-width: 700px;
  margin: 20px auto;
  font-family: sans-serif;
}
</style>