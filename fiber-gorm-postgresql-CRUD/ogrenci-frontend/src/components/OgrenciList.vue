<script setup>
defineProps({
  ogrenciler: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['sil', 'duzenle'])

function silmeyiOnayla(id) {
  if (confirm('Bu ogrenciyi silmek istediginizden emin misiniz?')) {
    emit('sil', id)
  }
}
</script>

<template>
  <div>
    <h2>Öğrenci Listesi</h2>
    <table>
      <thead>
        <tr>
          <th>Ad</th>
          <th>Soyad</th>
          <th>Numara</th>
          <th>Bölüm</th>
          <th>İşlemler</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="ogrenci in ogrenciler" :key="ogrenci.ID">
            <td>{{ ogrenci.ad }}</td>
            <td>{{ ogrenci.soyad }}</td>
            <td>{{ ogrenci.numara }}</td>
            <td>{{ ogrenci.bolum }}</td>
            <td>
                <button @click="emit('duzenle', ogrenci)">Düzenle</button>
                <button @click="silmeyiOnayla(ogrenci.ID)">Sil</button>
            </td>
            </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
}
th, td {
  border: 1px solid #ccc;
  padding: 8px;
  text-align: left;
}
</style>