import axios from "axios"

const API_URL = "http://localhost:3001/api/scan"

export default {
  tara(ip, cidr) {
    return axios.post(API_URL, { ip, cidr })
  }
}