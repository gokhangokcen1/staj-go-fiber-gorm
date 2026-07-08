<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  duzenlenecek: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['kaydet'])

const ad = ref('')
const soyad = ref('')
const numara = ref('')
const bolum = ref('')

watch(() => props.duzenlenecek, (yeniDeger) => {
  if (yeniDeger) {
    ad.value = yeniDeger.ad
    soyad.value = yeniDeger.soyad
    numara.value = yeniDeger.numara
    bolum.value = yeniDeger.bolum
  } else {
    formuTemizle()
  }
})

function formuTemizle() {
  ad.value = ''
  soyad.value = ''
  numara.value = ''
  bolum.value = ''
}

function gonder() {
  if (!ad.value || !numara.value) {
    alert('Ad ve Numara alanları zorunludur.')
    return
  }

  emit('kaydet', {
    ad: ad.value,
    soyad: soyad.value,
    numara: numara.value,
    bolum: bolum.value
  })

  formuTemizle()
}
</script>

<template>
  <form @submit.prevent="gonder">
    <h2>{{ duzenlenecek ? 'Öğrenci Güncelle' : 'Yeni Öğrenci Ekle' }}</h2>

    <input v-model="ad" placeholder="Ad" />
    <input v-model="soyad" placeholder="Soyad" />
    <input v-model="numara" placeholder="Numara" />
    <input v-model="bolum" placeholder="Bolum" />

    <button type="submit">{{ duzenlenecek ? 'Güncelle' : 'Ekle' }}</button>

  </form>
</template>

<style scoped>
form {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
}
input, button {
  padding: 8px;
}
</style>