<template lang="pug">
  v-container(fluid)
    v-row(v-if="orders" dense)
      v-col.mt-4
        h3 Order list
        v-data-table(
          :items="orders"
          :headers="headers"
          class="elevation-1"
          hide-default-footer
          disable-sort
        )
          template(v-slot:item.id="{ item }")
            router-link(:to="{name: 'OrderForm', params: {id: item.id}}") {{ item.id }}
          template(v-slot:item.status="{ item }")
            v-chip.text-capitalize(dark :color="item.status === 'paid' ? 'green' : ''") {{ item.status }}
          template(v-slot:item.total_for_payment="{ item }")
            .font-weight-bold {{ sumFromCents(item.total_for_payment).toFixed(2) }} $
          template(v-slot:item.action="{ item }")
            v-btn(v-if="item.status === 'pending'" dark color="orange darken-2" :to="{name: 'PaymentForm', params: {order_id: item.id}}") Pay
    v-row(v-else)
      v-col.d-flex.align-center.flex-column
        h2 You have no orders yet.
        span You can add products to your cart then checkout them.
        span All you orders will be here!
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator'
import api from '@/api/client'
import eventBus from '@/utils/event_bus'

@Component
export default class OrderList extends Vue {
  private orders = []
  private headers = [
    {text: 'Id', value: 'id'},
    {text: 'Status', value: 'status'},
    {text: 'Total for payment', value: 'total_for_payment'},
    {action: '', value: 'action', width: "100"}
  ]

  private created() {
    this.getOrders()
  }

  private sumFromCents(cents: number): number {
    return (cents / 100)
  }

  private getOrders() {
    api.get('/order').then((response) => {
      this.orders = response.data.data
    }).catch((err) => {
      eventBus.$emit('app_err', err)
    })
  }
}
</script>