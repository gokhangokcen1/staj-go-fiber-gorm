<script setup>
import { ref } from 'vue'
import SSLCheckService from '../services/SSLCheckService'

const website = ref('')
const port = ref(443)
const sonuc = ref(null)
const hataMesaji = ref('')
const yukleniyor = ref(false)

async function kontrolEt() {
  hataMesaji.value = ''
  sonuc.value = null
  yukleniyor.value = true

  try {
    const response = await SSLCheckService.sslcheck(website.value, Number(port.value))
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
    <h1>SSL Checker</h1>

    <form @submit.prevent="kontrolEt">
      <input v-model="website" placeholder="Örnek: google.com" />
      <input v-model="port" type="number" min="1" max="65535" placeholder="Port (ornek: 443)" />
      <button type="submit" :disabled="yukleniyor">
        {{ yukleniyor ? 'Kontrol Ediliyor...' : 'Kontrol Et' }}
      </button>
    </form>

    <p v-if="hataMesaji" class="hata">{{ hataMesaji }}</p>

    <div v-if="sonuc" class="sonuc">
      <h2>Genel Bilgi</h2>
      <table>
        <tbody>
          <tr><td>Resolves to</td><td>{{ sonuc.General.ResolvesTo }}</td></tr>
          <tr><td>Son geçerlilik</td><td>{{ sonuc.General.ExpirationDate }}</td></tr>
          <tr><td>Vendor imzalı</td><td>{{ sonuc.General.VendorSigned ? 'Evet' : 'Hayir' }}</td></tr>
          <tr><td>Hostname eşleşiyor</td><td>{{ sonuc.General.HostnameMatches ? 'Evet' : 'Hayir' }}</td></tr>
          <tr><td>Key length</td><td>{{ sonuc.General.KeyLength }}</td></tr>
          <tr><td>Server type</td><td>{{ sonuc.General.ServerType }}</td></tr>
          <tr><td>Revocation status</td><td>{{ sonuc.General.RevocationStatus }}</td></tr>
        </tbody>
      </table>

      <h2>Issued For</h2>
      <table>
        <tbody>
          <tr><td>Common Name</td><td>{{ sonuc.For.CommonName }}</td></tr>
          <tr><td>Organization</td><td>{{ sonuc.For.Organization }}</td></tr>
          <tr><td>Country</td><td>{{ sonuc.For.Country }}</td></tr>
          <tr><td>SAN</td><td>{{ sonuc.For.SAN.join(', ') }}</td></tr>
          <tr><td>Organization unit</td><td>{{ sonuc.For.OrganizationUnit }}</td></tr>
          <tr><td>State</td><td>{{ sonuc.For.State }}</td></tr>
          <tr><td>Locality</td><td>{{ sonuc.For.Locality }}</td></tr>
          <tr><td>Address</td><td>{{ sonuc.For.Address }}</td></tr>
        </tbody>
      </table>

      <h2>Issued By</h2>
      <table>
        <tbody>
          <tr><td>Common Name</td><td>{{ sonuc.By.CommonName }}</td></tr>
          <tr><td>Organization</td><td>{{ sonuc.By.Organization }}</td></tr>
          <tr><td>Country</td><td>{{ sonuc.By.Country }}</td></tr>
          <tr><td>Organization unit</td><td>{{ sonuc.By.OrganizationUnit }}</td></tr>
          <tr><td>State</td><td>{{ sonuc.By.State }}</td></tr>
          <tr><td>Locality</td><td>{{ sonuc.By.Locality }}</td></tr>
        </tbody>
      </table>

      <h2>Chain Details</h2>
      <table>
        <thead>
          <tr>
            <th></th>
            <th v-for="(cert, i) in sonuc.Chain.Certs" :key="i">Cert {{ i + 1 }}</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><strong>Issuer</strong></td>
            <td v-for="(cert, i) in sonuc.Chain.Certs" :key="i">{{ cert.Issuer }}</td>
          </tr>
          <tr>
            <td><strong>Common name</strong></td>
            <td v-for="(cert, i) in sonuc.Chain.Certs" :key="i">{{ cert.CommonName }}</td>
          </tr>
          <tr>
            <td><strong>Organization</strong></td>
            <td v-for="(cert, i) in sonuc.Chain.Certs" :key="i">{{ cert.Organization }}</td>
          </tr>
          <tr>
            <td><strong>Issued</strong></td>
            <td v-for="(cert, i) in sonuc.Chain.Certs" :key="i">{{ cert.Issued }}</td>
          </tr>
          <tr>
            <td><strong>Expires</strong></td>
            <td v-for="(cert, i) in sonuc.Chain.Certs" :key="i">{{ cert.Expires }}</td>
          </tr>
        </tbody>
      </table>
            <h2>Advanced</h2>
      <table>
        <tbody>
          <tr><td>Serial number</td><td>{{ sonuc.Chain.Certs[0].SerialNumber }}</td></tr>
          <tr><td>Signature algorithm</td><td>{{ sonuc.Chain.Certs[0].SignatureAlgorithm }}</td></tr>
          <tr><td>Fingerprint (SHA-1)</td><td>{{ sonuc.Chain.Certs[0].FingerprintSHA1 }}</td></tr>
          <tr><td>Fingerprint (MD5)</td><td>{{ sonuc.Chain.Certs[0].FingerprintMD5 }}</td></tr>
        </tbody>
      </table>
      <h2>Certificate</h2>
      <pre>{{ sonuc.Chain.Certs[0].PEM }}</pre>
    </div>
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
  margin-bottom: 24px;
}
td, th {
  border: 1px solid #ccc;
  padding: 8px;
  text-align: left;
}
</style>