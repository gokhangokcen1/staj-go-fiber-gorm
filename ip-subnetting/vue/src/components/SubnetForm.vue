<script setup>
import { ref } from 'vue'
import SubnetService from '../services/SubnetService'

const ip = ref('')
const cidr = ref(24)
const sonuc = ref(null)
const hataMesaji = ref('')

async function hesapla() {
  hataMesaji.value = ''
  sonuc.value = null

  try {
    const response = await SubnetService.hesapla(ip.value, Number(cidr.value))
    sonuc.value = response.data
  } catch (error) {
    if (error.response && error.response.data && error.response.data.error) {
      hataMesaji.value = error.response.data.error
    } else {
      hataMesaji.value = 'Hesaplama basarisiz, backend calisiyor mu kontrol et'
    }
  }
}
</script>

<template>
  <div>
    <h1>IP Subnetting Hesaplayici</h1>

    <form @submit.prevent="hesapla">
      <input v-model="ip" placeholder="Ornek: 192.168.1.10" />
      <input v-model="cidr" type="number" min="0" max="32" placeholder="CIDR (ornek: 24)" />
      <button type="submit">Hesapla</button>
    </form>

    <p v-if="hataMesaji" class="hata">{{ hataMesaji }}</p>

    <table v-if="sonuc">
      <tbody>
<tr>
    <td>Adres</td>
    <td>{{ sonuc.adres }}</td>
    <td>{{ sonuc.addressBinary }}</td>
</tr>

<tr>
    <td>Network Mask</td>
    <td>{{ sonuc.networkMask }}</td>
    <td>{{ sonuc.maskBinary }}</td>
</tr>

<tr>
    <td>Wildcard Mask</td>
    <td>{{ sonuc.wildcardMask }}</td>
    <td>{{ sonuc.wildcardMaskBinary }}</td>
</tr>

<tr>
    <td>Network Address</td>
    <td>{{ sonuc.networkAddress }}</td>
    <td>{{ sonuc.networkBinary }}</td>
</tr>

<tr>
    <td>Broadcast Address</td>
    <td>{{ sonuc.broadcastAddress }}</td>
    <td>{{ sonuc.broadcastBinary }}</td>
</tr>

<tr>
    <td>Hostmin</td>
    <td>{{ sonuc.hostmin }}</td>
    <td>{{ sonuc.hostMinBinary }}</td>
</tr>

<tr>
    <td>Hostmax</td>
    <td>{{ sonuc.hostmax }}</td>
    <td>{{ sonuc.hostMaxBinary }}</td>
</tr>

<tr>
    <td>Hosts/Net</td>
    <td>{{ sonuc.hostsPerNet }}</td>
    <td></td>
</tr>
<tr>
    <td>Class</td>
    <td>{{ sonuc.class }}</td>
</tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
form {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
input, button {
  padding: 8px;
}
.hata {
  color: red;
}
table {
  width: 100%;
  border-collapse: collapse;
}
td {
  border: 1px solid #ccc;
  padding: 8px;
}
</style>