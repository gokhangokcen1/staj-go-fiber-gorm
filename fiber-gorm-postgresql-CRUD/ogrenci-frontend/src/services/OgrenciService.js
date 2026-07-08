import axios from "axios"

const API_URL = "http://localhost:3000/api/ogrenci"

export default {
    getAll(){
        return axios.get(API_URL)
    },
    getById(id){
        return axios.get(`${API_URL}/${id}`)
    },
    create(ogrenci) {
        return axios.post(API_URL, ogrenci)
    },
    update(id, ogrenci) {
        return axios.put(`${API_URL}/${id}`, ogrenci)
    },
    delete(id) {
        return axios.delete(`${API_URL}/${id}`)
    }
}