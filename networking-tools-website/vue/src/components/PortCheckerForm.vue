<script setup>
import { ref } from 'vue'
import PortCheckerService from '../services/PortCheckerService'

const yayginPortlar = [
  { isim: 'FTP - 21', deger: 21 },
  { isim: 'SSH - 22', deger: 22 },
  { isim: 'SMTP - 25', deger: 25 },
  { isim: 'DNS - 53', deger: 53 },
  { isim: 'HTTP - 80', deger: 80 },
  { isim: 'POP3 - 110', deger: 110},
  { isim: 'IMAP - 143', deger: 143},
  { isim: 'IRC - 194', deger: 194},
  { isim: 'HTTPS - 443', deger: 443 },
  { isim: 'SMB - 445', deger: 445},
  { isim: 'SMTPS - 465', deger: 465},
  { isim: 'MySQL - 3306', deger: 3306 },
  { isim: 'REMOTE DESKTOP - 3389', deger: 3389},
  { isim: 'PostgreSQL - 5432', deger: 5432 },
]

const ip = ref('')
const port = ref(443)
const sonuc = ref(null)
const hataMesaji = ref('')
const yukleniyor = ref(false)

async function mevcutIPKullan() {
  try {
    const response = await PortCheckerService.mevcutIPGetir()
    ip.value = response.data.ip
  } catch (error) {
    hataMesaji.value = 'IP alinamadi, backend calisiyor mu kontrol et'
  }
}

async function kontrolEt() {
  hataMesaji.value = ''
  sonuc.value = null
  yukleniyor.value = true

  try {
    const response = await PortCheckerService.kontrolEt(ip.value, Number(port.value))
    sonuc.value = response.data
  } catch (error) {
    if (error.response && error.response.data && error.response.data.error) {
      hataMesaji.value = error.response.data.error
    } else {
      hataMesaji.value = 'Kontrol basarisiz, backend calisiyor mu kontrol et'
    }
  } finally {
    yukleniyor.value = false
  }
}
</script>

<template>
  <div>
    <h1>Port Checker</h1>

    <label>IP Adresi</label>
    <div class="ip-satiri">
      <input v-model="ip" placeholder="Örnek: 192.168.1.10" />
      <button @click="mevcutIPKullan">Mevcut IP'yi Kullan</button>
    </div>

    <label>Port Numarası</label>
    <div class="port-satiri">
      <input v-model="port" type="number" min="1" max="65535" />
      <select v-model="port">
        <option v-for="p in yayginPortlar" :key="p.deger" :value="p.deger">
          {{ p.isim }}
        </option>
      </select>
    </div>

    <button class="kontrol-butonu" @click="kontrolEt" :disabled="yukleniyor">
      {{ yukleniyor ? 'Kontrol Ediliyor...' : 'Tara' }}
    </button>

    <p v-if="hataMesaji" class="hata">{{ hataMesaji }}</p>

    <div v-if="sonuc" class="sonuc">
      Port {{ sonuc.port }} is
      <strong :class="sonuc.acik ? 'acik' : 'kapali'">
        {{ sonuc.acik ? 'open' : 'closed' }}
      </strong>.
    </div>
  </div>
</template>

<style scoped>
label {
  display: block;
  margin-top: 16px;
  font-weight: bold;
}
.ip-satiri, .port-satiri {
  display: flex;
  gap: 8px;
  margin-top: 4px;
}
input, select, button {
  padding: 8px;
}
.kontrol-butonu {
  margin-top: 24px;
  width: 100%;
  padding: 12px;
  background: #222;
  color: white;
  cursor: pointer;
}
.hata {
  color: red;
  margin-top: 16px;
}
.sonuc {
  margin-top: 24px;
  padding: 16px;
  background: #eee;
  font-style: italic;
}
.acik {
  color: green;
}
.kapali {
  color: red;
}
</style>