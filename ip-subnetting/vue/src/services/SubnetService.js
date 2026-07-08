import axios from "axios"

const API_URL = "http://localhost:3001/api/subnet"

export default {
  hesapla(ip, cidr) {
    return axios.post(API_URL, { ip, cidr })
  }
}
