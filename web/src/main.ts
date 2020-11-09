import Vue from 'vue'
import App from './App.vue'
import createRouter from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import './plugins/element_ui'
import './style/index.scss'

Vue.config.productionTip = false

new Vue({
  router: createRouter(),
  store,
  vuetify,
  render: h => h(App),
  async created() {
    (this.$router as any).start()
  }
}).$mount('#app')
