<template lang="pug">
  v-container(fluid)
    v-row(v-if="showView" dense)
      v-col(
        cols="12"
        v-if="cart.non_discount_set"
        v-for="(item, idx) in cart.non_discount_set"
        :key="'nd'+idx"
      )
        CartItem(:is-loading="loading" :cart-item="item" @delete-item="handleDeleteItem" :items="cart.items" @change-item-amount="handleChangeItemAmount")
      v-col(
        cols="12"
        v-if="cart.discount_sets"
        v-for="(set, idx) in cart.discount_sets"
        :key="'d'+idx"
      )
        template(
          v-if="Object.entries(set).length === 1"
        )
          CartItem(:is-loading="loading" :cart-item="set[0]" @delete-item="handleDeleteItem" :items="cart.items" @change-item-amount="handleChangeItemAmount")
        v-card(
          v-if="Object.entries(set).length > 1"
        )
          v-card-text
            v-row(dense).black--text.align-center
              v-col(cols="6")
                v-list
                  template(v-for="(item, idx) in set")
                    v-subheader.red--text.font-weight-bold(v-if="idx === 0") {{ item.discount_percent }}% set discount
                    CartSetItem(:key="idx" :is-loading="loading" :cart-item="item" @delete-item="handleDeleteItem" :items="cart.items" @change-item-amount="handleChangeItemAmount")
                    v-divider(v-if="idx < Object.entries(set).length - 1")
              v-spacer
              v-col.d-flex.justify-end.align-end.flex-column
                span.body-1 Total:
                span.body-1.text-decoration-line-through.grey--text.mr-1 {{ calcSetTotal(set).toFixed(2) }}$
                span.body-1.font-weight-bold {{ calcSetDiscountedTotal(set).toFixed(2) }}$
    v-row.mt-4(v-if="showView")
      v-col(cols="12" md="5").d-flex.align-center
        v-text-field.mr-4(v-model="couponCode" label="Coupon code" outlined dense hide-details)
        v-btn(dark color="orange" @click="applyCoupon") Apply
    v-row(v-if="showView && cart.coupons.length > 0")
      v-col
        h3 Coupons
        v-data-table.mt-4(
          :headers="couponsTableHeaders"
          :items="cart.coupons"
          class="elevation-1"
          hide-default-footer
          disable-sort
        )
          template(v-slot:item.expire="{ item }")
            span {{ item.expire | parseTimestamp }}
          template(v-slot:item.actions="{ item }")
            v-btn(dark color="red" @click="detachCoupon(item)") Remove
    v-row(v-if="showView")
      v-col
        v-btn(dark color="green" @click="checkout") Checkout
      v-spacer
      v-col.d-flex.flex-column.align-end(v-if="cart")
        span Total: {{ sumFromCents(cart.total).toFixed(2) }} $
        span Savings total: {{ sumFromCents(cart.total_saving).toFixed(2) }} $
        span.mt-4.font-weight-bold Total for payment {{ sumFromCents(cart.total_for_payment).toFixed(2) }} $
    v-row(v-else)
      v-col.d-flex.align-center.flex-column
        h2 No items in the cart
        span You can add products from the main page!
</template>

<script lang="ts">
import {Component, Mixins} from 'vue-property-decorator'
import api from '@/api/client'
import sortCartItems from '@/utils/order_object_keys'
import AppMixin from '@/mixins/AppMixin'
import CartItem from '@/views/components/CartItem.vue'
import CartSetItem from '@/views/components/CartSetItem.vue'
import _debounce from 'lodash/debounce'

@Component({
  components: {CartItem, CartSetItem}
})
export default class Cart extends Mixins(AppMixin) {
  private cart: any = null
  private couponCode = ''
  private loading = false

  private couponsTableHeaders = [
    {text: 'Code', value: 'code'},
    {text: 'Name', value: 'name'},
    {text: 'Expire At', value: 'expire'},
    {text: 'Percent', value: 'percent'},
    {text: '', value: 'actions', width: '100'}
  ]

  get showView(): boolean {
    return this.cart && Object.entries(this.cart.items).length > 0
  }

  private handleDeleteItem(item: any) {
    if (item) {
      this.setCartItemAmount(item.product.id, 0)
    }
  }

  private handleChangeItemAmount(obj: any) {
    if (obj && obj.item) {
      console.log('handleChangeItemAmount', obj)
      this.debouncedSetCartItemAmount(obj.item.product.id, obj.amount)
    }
  }

  private debouncedSetCartItemAmount = _debounce(this.setCartItemAmount, 100)

  private tt(id: string, amount: number) {
    console.log('TT', id, amount)
  }

  public created() {
    this.getCart()
  }
  
  private calcItemTotal(item: any): number {
    return (item.product.price / 100 * item.amount)
  }

  private calcItemDiscountedTotal(item: any): number {
    let total = item.product.price / 100 * item.amount
    total = total - (total * item.discount_percent / 100)
    return total
  }

  private calcSetTotal(set: any[]): number {
    let total = 0
    set.forEach((item: any) => {
      total += this.calcItemTotal(item)
    })
    return total
  }

  private calcSetDiscountedTotal(set: any[]): number {
    let total = 0
    set.forEach((item: any) => {
      total += this.calcItemDiscountedTotal(item)
    })
    return total
  }

  private applyCoupon() {
    api.post(`/cart/coupon`, {code: this.couponCode}).then((response) => {
      this.couponCode = ''
      this.handleGetCartResponse(response)
    }).catch((err) => {
      this.showError(err)
    })
  }

  private detachCoupon(item: any) {
    api.delete(`/cart/coupon/${item.code}`).then((response) => {
      this.handleGetCartResponse(response)
    }).catch((err) => {
      this.showError(err)
    })
  }

  private setCartItemAmount(productId: string, amount: number) {
    this.loading = true

    api.post('/cart/product', {
      product_id: productId,
      amount: amount
    }).then((response: any) => {
      this.handleGetCartResponse(response)
    }).catch((err) => {
      this.showError(err)
    }).then(() => {
      this.loading = false
    })
  }

  private getCart() {
    this.loading = true

    api.get('/cart').then((response) => {
      this.handleGetCartResponse(response)
    }).catch((err) => {
      this.showError(err)
    }).then(() => {
      this.loading = false
    })
  }

  private checkout() {
    api.put('/order').then((response) => {
      const order = response.data.data
      this.$router.push({name: 'OrderForm', params: {id: order.id, order: order}})
    }).catch((err) => {
      this.showError(err)
    })
  }

  private handleGetCartResponse(response: any) {
    this.cart = sortCartItems(response.data.data)
  }
}
</script>