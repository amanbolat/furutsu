<template lang="pug">
  v-card.w-100.elevation-3
    v-card-title
      | Registration Form
    v-card-text
      v-form(ref="form" v-model="valid")
        v-text-field(
          :disabled="isLoading"
          v-model="registrationForm.full_name"
          type="text"
          label="Full Name"
          name="full_name"
          required
          :rules="fullNameRule")
        v-text-field(
          :disabled="isLoading"
          v-model="registrationForm.username"
          type="text"
          label="Username"
          name="username"
          required
          :rules="usernameRules")
        v-text-field(
          :disabled="isLoading"
          v-model="registrationForm.password"
          type="password"
          label="Password"
          name="password"
          required
          :rules="passwordRules"
          @keyup.enter.native="handleRegister")
    v-card-actions
      v-spacer
      v-btn(@click="close") Cancel
      v-btn(
        dark
        color="indigo"
        @click.native.prevent="handleRegister"
        :loading="isLoading") Submit
</template>

<script lang="ts">
import {Component, Emit, Mixins, Prop} from 'vue-property-decorator'
import AppMixin from '@/mixins/AppMixin'
import api from '@/api/client'

@Component({
  name: 'RegistrationForm',
})
export default class RegistrationForm extends Mixins(AppMixin) {
  private isLoading = false
  private registrationForm = {
    username: '',
    password: '',
    full_name: ''
  }
  private valid = true

  @Prop([Function]) public successRegisterHandler!: () => {}

  @Emit()
  public close() {
    this.registrationForm.username = ''
    this.registrationForm.password = ''
  }

  get usernameRules(): any[] {
    return [
      (v: string) => !!v || 'Username is required',
      (v: string) => (!!v && v.length >= 3 && v.length <= 32) || 'Username should be between 3 and 32 characters',
  ]
  }

  get passwordRules(): any[] {
    return [
      (v: string) => !!v || 'Password is required',
      (v: string) => (!!v && v.length >= 3 && v.length <= 32) || 'Password should be between 3 and 32 characters',
    ]
  }

  get fullNameRule(): any[] {
    return [
      (v: string) => !!v || 'Full Name is required',
      (v: string) => (!!v && v.length >= 3 && v.length <= 100) || 'Full Name should be between 3 and 100 characters'
    ]
  }

  public handleRegister() {
    if ((this.$refs.form as any).validate()) {
      this.isLoading = true
      api.post('/auth/register', this.registrationForm).then(() => {
        this.isLoading = false
        this.close()
        this.successRegisterHandler()
        this.showMessage(`Welcome ${this.registrationForm.full_name}. Now you can use your username and password to sign in`)
      }).catch((error: any) => {
        this.isLoading = false
        this.showError(error)
      })
    } else {
      return
    }
  }
}
</script>
