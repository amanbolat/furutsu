<template lang="pug">
  v-container(fluid)
    v-row
      v-col(v-for="(item, idx) in products" :key="idx" cols="12" md="6")
        v-card.mx-auto(max-width='344')
          v-img(:src='`/img/${item.name}.jpg`' height='200px')
          v-card-title.text-capitalize
            | {{ item.name }}
          v-card-subtitle
            | Price: {{ item.price / 100 }} $
          v-card-actions
            v-menu(offset-y :close-on-content-click="false")
              template(v-slot:activator="{on, attrs}")
                v-btn(color='orange darken-2' text v-on="on" v-bind="attrs") Add to cart
              v-card
                v-card-text.d-flex.justify-center
                  a-input-number(v-model="addToCartAmount" :min="0" :step="1")
                v-card-actions
                  v-btn(@click="addProductToCart(item, addToCartAmount)") Add
            v-spacer
          v-expansion-panels
            v-expansion-panel
              v-expansion-panel-header
                | Description
              v-expansion-panel-content
                | {{ item.description }}

</template>

<script lang="ts">
import {Component, Mixins} from 'vue-property-decorator'
import api from '@/api/client'
import AppMixin from '@/mixins/AppMixin'

@Component
export default class ProductList extends Mixins(AppMixin) {
  private products: object[] = []
  private addToCartmenu = false
  private addToCartAmount = 0

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

  private addProductToCart(product: any, amount: number) {
    if (amount < 1) {
      return
    }
    api.put('/cart/product', {
      product_id: product.id,
      amount: amount
    }).then(() => {
      this.addToCartAmount = 0
      this.showMessage(`Added ${amount} ${product.name}s to the cart`)
    }).catch((err) => {
      this.showError(err)
    })
  }
}
</script>
