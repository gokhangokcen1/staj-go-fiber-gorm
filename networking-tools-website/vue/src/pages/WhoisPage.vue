<script setup>
import { ref } from 'vue'
import { lookupWhois } from '@/services/WhoisService'

const domain = ref('')
const report = ref(null)
const loading = ref(false)
const error = ref('')

async function handleLookup() {
  loading.value = true
  error.value = ''
  report.value = null
  try { report.value = (await lookupWhois(domain.value.trim())).data }
  catch (e) { error.value = e.response?.data?.error || e.message || 'WHOIS sorgusu basarisiz oldu.' }
  finally { loading.value = false }
}
function value(item) { return item || 'Bilgi yok' }
</script>

<template>
  <div class="whois-page page-card">
    <h2>WHOIS Sorgusu</h2>
    <p class="intro">Bir alan adi veya URL girin; kayit, kayit firmasi ve iletisim bilgilerini goruntuleyin.</p>
    <form @submit.prevent="handleLookup">
      <input v-model="domain" type="text" placeholder="tesla.com veya https://tesla.com" autocomplete="url" required />
      <button type="submit" :disabled="loading">{{ loading ? 'Sorgulaniyor...' : 'Sorgula' }}</button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
    <div v-if="report" class="whois-results">
      <section class="whois-section"><h3>Domain Information</h3><dl class="details-grid">
        <dt>Domain</dt><dd>{{ report.domain }}</dd><dt>Registered On</dt><dd>{{ value(report.registeredOn) }}</dd><dt>Expires On</dt><dd>{{ value(report.expiresOn) }}</dd><dt>Updated On</dt><dd>{{ value(report.updatedOn) }}</dd>
        <dt>Status</dt><dd><ul v-if="report.status?.length"><li v-for="status in report.status" :key="status">{{ status }}</li></ul><span v-else>Bilgi yok</span></dd>
        <dt>Name Servers</dt><dd><ul v-if="report.nameServers?.length"><li v-for="server in report.nameServers" :key="server">{{ server }}</li></ul><span v-else>Bilgi yok</span></dd>
      </dl></section>
      <section class="whois-section"><h3>Registrar Information</h3><dl class="details-grid">
        <dt>Registrar</dt><dd>{{ value(report.registrar?.name) }}</dd><dt>IANA ID</dt><dd>{{ value(report.registrar?.ianaId) }}</dd><dt>Email</dt><dd>{{ value(report.registrar?.email) }}</dd><dt>Abuse Email</dt><dd>{{ value(report.registrar?.abuseEmail) }}</dd><dt>Abuse Phone</dt><dd>{{ value(report.registrar?.abusePhone) }}</dd>
      </dl></section>
      <section class="whois-section"><h3>Registrant Contact</h3><dl class="details-grid">
        <dt>Name</dt><dd>{{ value(report.registrant?.name) }}</dd><dt>Organization</dt><dd>{{ value(report.registrant?.organization) }}</dd><dt>Street</dt><dd>{{ value(report.registrant?.street) }}</dd><dt>City</dt><dd>{{ value(report.registrant?.city) }}</dd><dt>State</dt><dd>{{ value(report.registrant?.state) }}</dd><dt>Postal Code</dt><dd>{{ value(report.registrant?.postalCode) }}</dd><dt>Country</dt><dd>{{ value(report.registrant?.country) }}</dd><dt>Phone</dt><dd>{{ value(report.registrant?.phone) }}</dd><dt>Fax</dt><dd>{{ value(report.registrant?.fax) }}</dd><dt>Email</dt><dd>{{ value(report.registrant?.email) }}</dd>
      </dl></section>
      <section class="whois-section"><h3>Technical Contact</h3><dl class="details-grid"><dt>Name</dt><dd>{{ value(report.technical?.name) }}</dd><dt>Phone</dt><dd>{{ value(report.technical?.phone) }}</dd><dt>Email</dt><dd>{{ value(report.technical?.email) }}</dd></dl></section>
    </div>
  </div>
</template>

<style scoped>
.intro { margin: -8px 0 20px; color: #a9bed2; }.whois-results { display: grid; gap: 16px; }.whois-section { padding: 21px; border: 1px solid rgba(173, 206, 231, .14); border-radius: 16px; background: rgba(7, 22, 38, .56); }.whois-section h3 { margin: 0 0 16px; color: #9df1dc; font-size: 1rem; }.details-grid { display: grid; grid-template-columns: minmax(130px, 190px) minmax(0, 1fr); margin: 0; }dt, dd { margin: 0; padding: 10px 0; border-top: 1px solid rgba(173, 206, 231, .1); overflow-wrap: anywhere; }dt { color: #a9bed2; font-size: .84rem; font-weight: 700; }dd { color: #e3eff9; font-family: 'DM Mono', ui-monospace, monospace; font-size: .85rem; }dd ul { display: grid; gap: 5px; margin: 0; padding: 0; list-style: none; }@media (max-width: 560px) { .details-grid { grid-template-columns: 1fr; } dt { padding-bottom: 3px; } dd { padding-top: 3px; border-top: 0; } }
</style>
