<template lang="pug">
  v-container(fluid)
    v-row(v-if="showView" dense)
      v-col(cols="12" md="6")
        h3 Order
        span {{ order.id }}
      v-spacer
      v-col.d-flex.justify-end
        v-chip(:color="order.status === 'paid' ? 'green' : ''" large).text-capitalize {{ order.status }}
    v-row(v-if="showView" dense).mt-4
      v-col
        v-data-table(
          :items="order.items"
          :headers="headers"
          class="elevation-1"
          hide-default-footer
          disable-sort
        )
          template(v-slot:item.price="{ item }")
            .font-weight-bold {{ sumFromCents(item.price).toFixed(2) }} $
    v-row(v-if="showView" dense).mt-4
      v-col
        v-btn(dark color="orange darken-2" @click="payForOrder" v-if="order.status === 'pending'") Pay
      v-spacer
      v-col.d-flex.flex-column.align-end
        span Total: {{ sumFromCents(order.total).toFixed(2) }} $
        span Savings total: {{ sumFromCents(order.savings).toFixed(2) }} $
        span.mt-4.font-weight-bold Total for payment {{ sumFromCents(order.total_for_payment).toFixed(2) }} $
</template>

<script lang="ts">
import {Component, Mixins} from 'vue-property-decorator'
import api from '@/api/client'
import AppMixin from '@/mixins/AppMixin'

@Component
export default class OrderForm extends Mixins(AppMixin)
{
  private order: any = {}

  private headers = [
    {text: 'Product', value: 'product_name'},
    {text: 'Amount', value: 'amount'},
    {text: 'Price', value: 'price'}
  ]

  private created() {
    this.getOrder()
  }

  get showView(): boolean {
    return this.order && Object.keys(this.order).length > 0
  }

  private payForOrder() {
    this.$router.push({name: 'PaymentForm', params: {order_id: this.order.id}})
  }

  private getOrder() {
    const idParam = this.$route.params.id
    const orderParam = this.$route.params.order

    if (orderParam) {
      this.order = orderParam
      return
    }

    api.get(`/order/${idParam}`).then((response) => {
      this.order = response.data.data
    }).catch((err) => {
     this.showError(err)
    })
  }
}
</script>