<template lang="pug">
  v-card(v-if="item && item.product")
    v-btn(icon @click="deleteItem(item)").delete-item-btn
      v-icon(color="red") mdi-trash-can-outline
    v-card-text
      v-row(dense).black--text.align-center
        v-col(cols="12" md="5")
          v-list-item
            v-list-item-avatar
              v-img(:src="`/img/${item.product.name}.jpg`")
            v-list-item-content
              v-list-item-title
                .text-capitalize.item-name {{ item.product.name }}
              v-list-item-subtitle
                span.black--text Price: {{ sumFromCents(item.product.price).toFixed(2)  }}$
                v-edit-dialog.mt-2(
                  :return-value.sync="item"
                  large
                  @open="openAmountDialog"
                  @save="saveAmountDialog"
                  @close="closeAmountDialog"
                  @cancel="closeAmountDialog")
                  .black--text.font-weight-bold.mr-2.amount-text.d-flex Amount: {{ item.amount }}
                    v-icon.ml-1(small) mdi-square-edit-outline
                  template(v-slot:input)
                    NumberInput(:disabled="loading" :value="item.amount" :min="1" :step="1" @change="handleAmountChange(item, $event)" inline controls size="small")
              v-list-item-subtitle(v-if="item.discount_percent")
                .red--text.font-weight-bold Discount: {{ item.discount_percent }}%
        v-spacer
        v-col.d-flex.justify-end.align-end.flex-column
          span.body-1 Total:
          span.body-1.text-decoration-line-through.grey--text.mr-1(v-if="cartItem.discount_percent") {{ calcItemTotal(cartItem).toFixed(2) }}$
          span.body-1.font-weight-bold {{ calcItemDiscountedTotal(cartItem).toFixed(2) }}$
</template>

<script lang="ts">
import {Component, Emit, Mixins, Prop, PropSync, Watch} from 'vue-property-decorator'
import _isNumber from 'lodash/isNumber'
import AppMixin from '@/mixins/AppMixin'

@Component({
  name: 'CartItem',
})
export default class CartItem extends Mixins(AppMixin) {
  @Prop(Object) public cartItem!: any
  @PropSync('items', {type: Object}) public cartItems!: any[]
  @PropSync('isLoading', {type: Boolean}) public loading!: any[]
  public item: any = {}
  public amountSaved = false

  @Watch('cartItem')
  onCartItemChange(val: any) {
    if (!val) {
      return
    }
    this.item = JSON.parse(JSON.stringify(val))
  }

  public created() {
    if (!this.cartItem) {
      return
    }
    this.item = JSON.parse(JSON.stringify(this.cartItem))
  }

  @Emit()
  private deleteItem(item: any) {
    return item
  }

  @Emit()
  public changeItemAmount(item: any, amount: number) {
    return {item, amount}
  }

  public handleAmountChange(item: any, amount: number) {
    if (!_isNumber(amount)) {
      return
    }
    if (amount < 1) {
      return
    }
    item.amount = amount
  }

  private openAmountDialog() {
    this.amountSaved = false
  }

  private closeAmountDialog() {
    if (this.amountSaved) {
      return
    }
    this.item.amount = this.cartItem.amount
  }

  private saveAmountDialog() {
    if (this.item.amount === this.cartItem.amount) {
      return
    }
    const diff = this.cartItems[this.item.product.id].amount + (this.item.amount - this.cartItem.amount)
    this.changeItemAmount(this.item, diff)
    this.amountSaved = true
  }

  private calcItemTotal(item: any): number {
    return (item.product.price / 100 * item.amount)
  }

  private calcItemDiscountedTotal(item: any): number {
    let total = item.product.price / 100 * item.amount
    total = total - (total * item.discount_percent / 100)
    return total
  }
}
</script>

<style lang="scss">
.item-name {
  font-size: 1.3em;
}

.v-small-dialog__activator {
  padding: 2px 0;
}

.amount-text {
  font-size: 1.1em;

  &:hover {
    border: thin dotted;
  }
}

.v-small-dialog__menu-content {
  padding-top: 10px;

  .v-small-dialog__content {
    display: flex;
    justify-content: center;
  }

  .v-small-dialog__actions {
    display: flex;
    justify-content: center;
  }
}

.delete-item-btn {
  position: absolute !important;
  bottom: 2px;
  right: 2px;
}
</style>