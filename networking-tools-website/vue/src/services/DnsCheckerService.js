import axios from 'axios'

const API_URL = 'http://localhost:3001/api'

export function checkDns(domain, recordType) {
  return axios.post(`${API_URL}/dnscheck`, { domain, recordType })
}
