<template lang="pug">
  v-app#inspire
    v-app-bar(app='' color='white' flat='')
      v-container.py-0.fill-height
        router-link(to="/")
          v-img(src="/logo.png" max-width="70")
        v-spacer
        v-btn(icon color="orange darken-2" to="/cart")
          v-icon mdi-cart
        v-btn(icon color="orange darken-2" to="/order").mr-4
          v-icon mdi-shopping
        v-dialog(v-if="!isAuthenticated" v-model="authDialog" width="500" persistent)
          template(v-slot:activator="{on, attrs}")
            v-btn(v-if="!isAuthenticated" v-on="on" v-bind="attrs") Login
          LoginForm(:success-login-handler="successLoginHandler" @close="closeLoginFormHandler")
        v-btn(v-if="isAuthenticated" @click="logout") Logout
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
import isFunc from 'lodash/isFunction'

@Component({
  name: 'App',
  components: {LoginForm}
})
export default class App extends Vue {
  private authDialog = false

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
    this.authDialog = false
  }

  public closeLoginFormHandler() {
    this.afterLoginRefused.forEach((f: () => void) => {
      if (isFunc(f)) {
        f()
      }
    })
    this.afterLoginRefused = []
    this.authDialog = false
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
    this.$set(this, 'authDialog', true)
    this.beforeSuccessLogin.push(handler)

    this.afterLoginRefused.push(() => {
      this.$router.push('/').catch(() => {
        console.log('handled auth prompt with error')
      })
    })
  }

  private handleUnauthorizedRequest() {
    this.$notify({
      title: 'Unauthorized',
      message: 'Please sign in'
    })
    this.$set(this, 'authDialog', true)

    this.beforeSuccessLogin.push(() =>{
      this.$router.go(0)
    })

    this.afterLoginRefused.push(() => {
      this.$router.push('/')
    })
  }

  private handleAppError(error: AppError) {
    console.log('APP ERR')

    let err = {} as AppError
    if (!error.message) {
      err.message = 'Unknown error'
    } else {
      err = error
    }
    this.$notify({
      title: err.message,
      message: err.hint,
      type: 'error',
    })
  }
}
</script>
