<template lang="pug">
  v-container(fluid)
    v-row
      v-col(v-for="(item, idx) in products" :key="idx" cols="12" md="4")
        v-card.mx-auto(max-width='344')
          v-img(:src='`/img/${item.Name}.jpg`' height='200px')
          v-card-title
            | {{ item.Name }}
          v-card-subtitle
            | Price: {{ item.Price / 100 }} $
          v-card-actions
            v-btn(color='orange lighten-2' text='')
              | Add to cart
            v-spacer
          v-expansion-panels
            v-expansion-panel
              v-expansion-panel-header
                | Description
              v-expansion-panel-content
                | {{ item.Description }}

</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator'
import api from '@/api/client'

@Component
export default class ProductList extends Vue {
  private products: object[] = []

  public created() {
    this.getProductList()
  }

  private getProductList() {
    api.get('/product').then((response) => {
      console.log(response)
      this.products = response.data.data
    }).catch((err) => {
      console.log('get products error', err)
    })
  }
}
</script>
