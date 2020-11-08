import Vue from 'vue'
import Vuex from 'vuex'
import user from './user'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

const pathsToPersist = [
  'user.user',
  'user.authToken',
]

export default new Vuex.Store({
  plugins: [
    createPersistedState({
      key: 'clientServiceStore',
      paths: [...pathsToPersist]
    }),
  ],
  modules: {
    user
  }
})



