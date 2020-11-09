<template lang="pug">
  v-app#inspire
    v-app-bar(app='' color='white' flat='')
      v-container.py-0.fill-height
        router-link(to="/") Home
        v-spacer
        v-btn(icon color="indigo" to="/cart")
          v-icon mdi-cart
        v-btn(icon color="indigo").mr-4
          v-icon mdi-shopping
        v-btn(v-if="isAuthenticated" @click="logout") Logout
        v-dialog(v-else v-model="authDialog" width="500")
          template(v-slot:activator="{on, attrs}")
            v-btn(v-on="on" v-bind="attrs") Login
          LoginForm
    v-main.grey.lighten-3
      v-container
        v-row
          v-col
            v-sheet.pa-5(min-height='70vh' rounded='lg')
              router-view


</template>


<script lang="ts">
import Vue from 'vue'
import eventBus from '@/utils/event_bus'
import {Component} from 'vue-property-decorator'
import {AppError} from '@/api/base'
import store from './store'
import LoginForm from './views/components/login_form.vue'

@Component({
  name: 'App',
  components: {LoginForm}
})
export default class App extends Vue {
  private authDialog = false

  get isAuthenticated() {
    return store.getters.isAuthenticated
  }

  public created() {
    eventBus.$on('app_error', this.handleAppError)
    eventBus.$on('unauthorized_request', this.handleUnauthorizedRequest)
  }

  public beforeDestroy() {
    eventBus.$off('app_error')
    eventBus.$off('unauthorized_request')
  }

  private logout() {
    store.dispatch('Logout')
  }

  private handleUnauthorizedRequest(error: any) {
    this.$notify({
      title: 'Unauthorized',
      message: 'Please sign in'
    })
    if (this.$router.currentRoute.path !== '/') {
      this.$router.push('/')
    }
  }

  private handleAppError(error: AppError) {
    this.$notify({
      title: error.message,
      message: error.hint,
      type: 'error',
      duration: 4000
    })
  }
}
</script>
