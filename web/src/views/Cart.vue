<template lang="pug">
  v-container(fluid)
    v-row(v-if="showView" dense)
      v-col.mt-4
        h3 Products without discounts
      v-col(
        cols="12"
        v-if="cart.non_discount_set"
        v-for="(item, idx) in cart.non_discount_set"
        :key="'nd'+idx"
      )
        v-card
          v-card-text
            v-row(dense).black--text.align-center
              v-col(cols="12" md="5")
                v-list-item
                  v-list-item-avatar
                    v-img(:src="`/img/${item.product.name}.jpg`")
                  v-list-item-content
                    v-list-item-title
                      .text-capitalize {{ item.product.name }}
                    v-list-item-subtitle
                      | Price: {{ sumFromCents(item.product.price).toFixed(2)  }}$
                      el-input-number.ml-5(:value="item.amount" :min="0" :step="1" size="small" @change="handleChangeItemAmount(item, $event)")
              v-spacer
              //v-col.d-flex

              v-col.d-flex.justify-end.align-end.flex-column
                span.body-1 Total:
                span.body-1.font-weight-bold  {{ calcItemTotal(item).toFixed(2) }}$
      v-col.mt-4
        h3 Products with discounts
      v-col(
        cols="12"
        v-if="cart.discount_sets"
        v-for="(set, idx) in cart.discount_sets"
        :key="'d'+idx"
      )
        v-card(
          v-if="Object.entries(set).length === 1"
          :_set="item = set[0]"
        )
          v-card-text
            v-row(dense).black--text.align-center
              v-col(cols="12" md="5")
                v-list-item
                  v-list-item-avatar
                    v-img(:src="`/img/${item.product.name}.jpg`")
                  v-list-item-content
                    v-list-item-title
                      .text-capitalize {{ item.product.name }}
                    v-list-item-subtitle
                      | Price: {{ sumFromCents(item.product.price).toFixed(2)  }}$
                    v-list-item-subtitle
                      .red--text.font-weight-bold Discount: {{ item.discount_percent }}%
              v-spacer
              v-col.d-flex.justify-end.align-end.flex-column
                span.body-1 Total:
                span.body-1.text-decoration-line-through.grey--text.mr-1 {{ calcItemTotal(item).toFixed(2) }}$
                span.body-1.font-weight-bold {{ calcItemDiscountedTotal(item).toFixed(2) }}$
        v-card(
          v-if="Object.entries(set).length > 1"
        )
          v-card-text
            v-row(dense).black--text.align-center
              v-col(cols="6")
                v-list
                  template(v-for="(item, idx) in set")
                    v-subheader.red--text.font-weight-bold(v-if="idx === 0") {{ item.discount_percent }}% set discount
                    v-list-item(:key="idx")
                      v-list-item-avatar
                        v-img(:src="`/img/${item.product.name}.jpg`")
                      v-list-item-content
                        v-list-item-subtitle.text-capitalize.black--text {{ item.product.name }}
                        v-list-item-subtitle Price: {{ sumFromCents(item.product.price).toFixed(2) }}$
                          el-input-number.ml-5(:value="item.amount" :min="0" :step="1" size="small" @change="handleChangeItemAmount(item, $event)")
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
import {Component, Mixins, Vue} from 'vue-property-decorator'
import api from '@/api/client'
import sortCartItems from '@/utils/order_object_keys'
import AppMixin from '@/mixins/AppMixin'

@Component
export default class Cart extends Mixins(AppMixin) {
  private cart: any = null
  private couponCode = ''

  private couponsTableHeaders = [
    {text: 'Code', value: 'code'},
    {text: 'Name', value: 'value'},
    {text: 'Percent', value: 'percent'},
    {text: '', value: 'actions', width: '100'}
  ]

  get showView(): boolean {
    return this.cart && Object.entries(this.cart.items).length > 0
  }

  private handleChangeItemAmount(item: any, newAmount: number) {
    console.log(item)
    console.log(newAmount)
    const amountDiff = newAmount - item.amount
    const setAmount = this.cart.items[item.product.id].amount + amountDiff
    item.amount = newAmount

    this.setCartItemAmount(item.product.id, setAmount)
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
    api.post('/cart/product', {
      product_id: productId,
      amount: amount
    }).then((response: any) => {
      this.handleGetCartResponse(response)
    }).catch((err) => {
      this.showError(err)
    })
  }

  private getCart() {
    api.get('/cart').then((response) => {
      this.handleGetCartResponse(response)
    }).catch((err) => {
      this.showError(err)
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