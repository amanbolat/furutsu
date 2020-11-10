<template lang="pug">
  v-app
    v-app-bar(app='' color='white' flat='')
      v-container.py-0.fill-height
        router-link(to="/")
          v-img(src="/logo.png" max-width="70")
        v-spacer
        v-btn(icon color="orange darken-2" to="/cart")
          v-icon mdi-cart
        v-btn(icon color="orange darken-2" to="/order").mr-4
          v-icon mdi-shopping
        v-dialog(v-if="!isAuthenticated" v-model="loginDialog" width="500" persistent)
          template(v-slot:activator="{on, attrs}")
            v-btn(v-if="!isAuthenticated" v-on="on" v-bind="attrs" color="indigo" dark) Login
          LoginForm(:success-login-handler="successLoginHandler" @close="closeLoginFormHandler")
        v-dialog(v-if="!isAuthenticated" v-model="registerDialog" width="500" persistent)
          template(v-slot:activator="{on, attrs}")
            v-btn.ml-3(v-if="!isAuthenticated" v-on="on" v-bind="attrs") Sign up
          RegistrationForm(:success-register-handler="successRegistrationHandler" @close="closeRegistrationFormHandler")
        v-btn(v-if="isAuthenticated" @click="logout") Logout
    v-main.grey.lighten-3
      v-container
        v-row
          v-col
            v-sheet.pa-5(min-height='70vh' rounded='lg')
              transition(name="fade" mode="out-in")
                router-view
        notifications(group="notify" position="top left")
</template>


<script lang="ts">
import eventBus from '@/utils/event_bus'
import {Component, Mixins} from 'vue-property-decorator'
import {AppError} from '@/api/base'
import store from './store'
import LoginForm from './views/components/LoginForm.vue'
import RegistrationForm from './views/components/RegistrationForm.vue'
import isFunc from 'lodash/isFunction'
import AppMixin from '@/mixins/AppMixin'

@Component({
  name: 'App',
  components: {LoginForm, RegistrationForm}
})
export default class App extends Mixins(AppMixin) {
  private loginDialog = false
  private registerDialog = false

  get isAuthenticated() {
    return store.getters.isAuthenticated
  }

  private beforeSuccessLogin: any[] = []
  private afterLoginRefused: any[] = []

  public successLoginHandler() {
    this.beforeSuccessLogin.forEach((f: () => void) => {
      if (isFunc(f)) {
        f()
      }
    })
    this.beforeSuccessLogin = []
    this.afterLoginRefused = []
    this.loginDialog = false
  }

  public closeLoginFormHandler() {
    this.afterLoginRefused.forEach((f: () => void) => {
      if (isFunc(f)) {
        f()
      }
    })
    this.afterLoginRefused = []
    this.loginDialog = false
  }

  private successRegistrationHandler() {
    this.registerDialog = false
  }

  private closeRegistrationFormHandler() {
    this.registerDialog = false
  }

  public created() {
    eventBus.$on('app_error', this.handleAppError)
    eventBus.$on('prompt_auth', this.handlePromptAuth)
    eventBus.$on('unauthorized_request', this.handleUnauthorizedRequest)
  }

  public beforeDestroy() {
    eventBus.$off('app_error')
    eventBus.$off('prompt_auth')
    eventBus.$off('unauthorized_request')
  }

  private logout() {
    store.dispatch('Logout')
    this.$router.push('/').catch(() => {
      console.log('logged out with error')
    })
  }

  private handlePromptAuth(handler: () => void) {
    this.$set(this, 'loginDialog', true)
    this.beforeSuccessLogin.push(handler)

    this.afterLoginRefused.push(() => {
      this.$router.push('/').catch(() => {
        console.log('handled auth prompt with error')
      })
    })
  }

  private handleUnauthorizedRequest(wasAuthorized: boolean) {
    let title = 'Not logged in'
    if (wasAuthorized) {
      title = 'Your session has expired'
    }
    this.$notify({
      group: 'notify',
      title: title,
      type: 'warn',
      text: 'Please sign in',
      duration: 2000
    })
    this.$set(this, 'loginDialog', true)

    this.beforeSuccessLogin.push(() =>{
      this.$router.go(0)
    })

    this.afterLoginRefused.push(() => {
      this.$router.push('/').catch(() => {
        console.log('handled login refuse with error')
      })
    })
  }

  private handleAppError(error: AppError) {
    this.showError(error)
  }
}
</script>
