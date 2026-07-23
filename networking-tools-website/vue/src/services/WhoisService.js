import axios from 'axios'

const API_URL = 'http://localhost:3001/api'

export function lookupWhois(domain) {
  return axios.post(`${API_URL}/whois`, { domain })
}
