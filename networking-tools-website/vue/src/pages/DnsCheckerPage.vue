<script setup>
import { ref } from 'vue'
import { checkDns } from '@/services/DnsCheckerService'
import DnsCheckerForm from '@/components/DnsCheckerForm.vue'

const recordTypes = ['A', 'AAAA', 'CNAME', 'MX', 'NS', 'PTR', 'SRV', 'SOA', 'TXT', 'CAA', 'DS', 'DNSKEY']
const domain = ref('')
const recordType = ref('A')
const results = ref([])
const checkedDomain = ref('')
const loading = ref(false)
const error = ref('')

async function handleCheck() {
  loading.value = true
  error.value = ''
  results.value = []
  try {
    const response = await checkDns(domain.value.trim(), recordType.value)
    checkedDomain.value = response.data.domain
    results.value = response.data.results || []
  } catch (e) {
    error.value = e.response?.data?.error || e.message || 'Kontrol başarısız. Backend çalışıyor mu kontrol edin.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="dns-checker-page page-card">
    <h2>DNS Checker</h2>
    <p class="intro">Bir kayit turu secin; genel DNS sunucularinin cevabi haritada ve hemen altinda listelenir.</p>

    <form class="dns-form" @submit.prevent="handleCheck">
      <input v-model="domain" type="text" placeholder="tesla.com" required />
      <select v-model="recordType" aria-label="DNS kayıt turu">
        <option v-for="type in recordTypes" :key="type" :value="type">{{ type }}</option>
      </select>
      <button type="submit" :disabled="loading">
        {{ loading ? 'Kontrol ediliyor...' : 'Kontrol Et' }}
      </button>
    </form>

    <p v-if="error" class="error">{{ error }}</p>
    <DnsCheckerForm :results="results" :record-type="recordType" :domain="checkedDomain" />
  </div>
</template>

<style scoped>
.intro { color: #a9bed2; margin: -6px 0 18px; }
.dns-form { display: flex; gap: 8px; margin-bottom: 16px; }
.dns-form input { flex: 1 1 280px; }
.dns-form select { min-width: 100px; }
.error { color: #ff8f8f; }
</style>
