<template lang="pug">
  v-card.w-100.elevation-3
    v-card-title
      | Login form
    v-card-text
      v-form(ref="form" v-model="valid")
        v-text-field(
          :disabled="isLoading"
          v-model="loginForm.username"
          type="text"
          label="Username"
          name="username"
          required
          :rules="usernameRules")
        v-text-field(
          :disabled="isLoading"
          v-model="loginForm.password"
          type="password"
          label="Password"
          name="password"
          required
          :rules="passwordRules"
          @keyup.enter.native="handleLogin")
    v-card-actions
      v-spacer
      v-btn(
        dark
        color="indigo"
        @click.native.prevent="handleLogin"
        :loading="isLoading") Login
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator'
import eventBus from '../../utils/event_bus'

@Component({
  name: 'LoginForm',
})
export default class Login extends Vue {
  private isLoading = false
  private loginForm = {
    username: '',
    password: ''
  }
  private valid = true

  get usernameRules(): any[] {
    return [
      (v: any) => !!v || 'Username is required'
    ]
  }

  get passwordRules(): any[] {
    return [
      (v: any) => !!v || 'Password is required'
    ]
  }

  public handleLogin() {
    if ((this.$refs.form as any).validate()) {
      this.isLoading = true
      this.$store.dispatch('Login', this.loginForm).then(() => {
        this.isLoading = false
        this.$router.push({path: '/'}
        ).catch((err) => {
          return
        })
      }).catch((error: any) => {
        this.isLoading = false
        eventBus.$emit('app_error', error)
      })
    } else {
      return
    }
  }
}
</script>
