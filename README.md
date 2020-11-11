# Furutsu - an online shop for selling fruits
![Logo](https://github.com/amanbolat/furutsu/raw/master/web/public/logo.png)

## About
This is the project developed by me as a homework assessment for the technical interview.

## How to run
###Option 1
- `git clone https://github.com/amanbolat/furutsu`
- `make dc.run`
- `make init.data`
- Open browser and go to `localhost:8080`
- Register new user
- Sign in

A few coupons will be inserted after `init.data`:
- orange333
- orange777
- orange888

### Option 2
- Start the postgresql database using docker or natively
- Run the server`PORT=9033 DB_CONN_STRING='postgres://postgres:postgres@127.0.0.1:5432/furutsu?sslmode=disable' MIGRATES_DIR=<path to migrates folder>  go run server/cmd/main.go`
- Serve the SPA `cd web && npm run serve`

ATTENTION: it was tested only on MacOS. If you use Linux or Windows machine you might have to do some tweaks.

## Task
Develop an online e-commerce store for selling fruits, which contains the following features:

- Simple sign-up and login form.
- Browse the following products
    1. Apples
    2. Bananas
    3. Pears
    4. Oranges
- Add items to your cart
- Adjust quantity.
- Delete items from the cart.
- Apply coupons.
- Checkout your cart.

Mocked purchase (a payment gateway is not required, but a route must exist in the backend validating the payment).
An address does not need to be entered.

### Requirements
- Architecture diagrams.
- Single-page frontend app (cannot use an existing online stores such as Prestashop).
- Backend RESTful web service written in GoLang.
- Users must be able to return to their cart after closing the browser, and see the previous items that were added.

### Cart Rules
- If 7 or more apples are added to the cart, a 10% discount is applied to all apples.
- For each set of 4 pears and 2 bananas, a 30% discount is applied, to each set.
- These sets must be added to their own cart item entry.
- If pears or bananas already exist in the cart, this discount must be recalculated when new pears or bananas are added.
- A coupon code can be used to get a 30% discount on oranges, if applied to the cart, otherwise oranges are full price.
- A coupon can only be applied once.
- It has a configurable expiry timeout (10 seconds for testing purposes) once generated.

### The following totals must be shown:
Total price.
Total savings.


## TODO
Like every project there are no limits in improvement. This is the list of features and improvements that I would add

### Tests
- [ ] Add e2e  tests.
- [ ] Add CI/CD workflow for test automation.
- [ ] Add some edge cases to the unit tests. Such as discounts for the same product.

### Security. Authorization & Authentication
- [ ] Implement password `service` to check the strength of the password
- [ ] Encrypt passwords before saving to the Database. Plain text is a BIG NO.
- [ ] The storage for JWT sign keys might be added for every user using Redis. Thus, we don't use sessions, yet
it is possible to renew sign keys in order to expire some JWT tokens. For example, in the situation when user changed 
the password.

### Metrics, logging
- [ ] Add something like ELK for logging.
- [ ] Add `Prometheus` for metrics.
- [ ] Add tracing if the monolith application will be divided into small microservices.

### CI/CD
- [ ] Add ci/cd pipeline to automate all the deployments

### Database
- [ ] Some application logic might be moved to the Database
- [ ] Add backups
- [ ] Add audit
- [ ] Add row level policies

### Code
- [ ] Clean it
- [ ] Put frontend part of the project to separate repo 

## Questions and Suggestions

1. The requirement `sets must be added to their own cart item entry` leaves a lot of question and without 
further discussion of that makes it difficult to add some new features or change the business logic. 

Questions:
- What should be done to the amount of products which doesn't fit that discount set? For instance, we have 10 
pears and 5 bananas, then we get only two sets comprised of `8 pears and 4 bananas` in total and `2 pears + 1 banana`
with no discount applied in two more entries.
- What if customer added pears and bananas, but wants to checkout only the pears without deleting bananas
from the cart? 

Suggestion: 
From the UI perspective adding the items left outside of the discount set to different entry is possible, but not 
recommended. Therefore, I would like to suggest create only one cart item for the product and just show, that part
of its amount might get some discounts.
There are some examples of good implementation of similar logic. Chinese e-commerce giants 
such as `Taobao` and `JD.com` have more complex discount system and leave the choice to customers. They 
provide `% discounts`, `set discounts`, `minimum order discounts` and customers may choose them by themselves.

2. Adding the set of items has also affected the functionality of quantity adjustment. If we add or remove items from the
set it changes the quantity of that particular product in the cart, BUT what about removing the item? Should we remove all
the items of this kind or remove only the exact quantity of the product in the set?
Currently, when user presses the DELETE button all the items of this kind will be removed from the cart.

Suggestion: reconsider the idea of the set. See question #1.   

3. Given cart rules about discounts would seem to be a little vague. 

Questions: 
- Discounts intersection. What if two or more discounts are applicable to the same set of items?
- What if the coupon intersects with the discount? 

My thoughts:
As I mentioned before giving the choice to the customer is the best option. However, it might depend on the real world
business requirements. Anyway I would be grateful to have a chance to discuss that question, just for curiosity.

4. Logic of the coupon is not very clear. It seems, that coupons should be generated somehow and everyone who has a code
can use them. Also, it is claimed to be applied to the cart, but not the item cart which seems to me as a bad logic. Indeed, 
the implementation of such a logic simpler, but it might create a few problems in long term. In my opinion, we could attach
coupons to the cart items, thus the discount and coupon discount would be used together or, again, as I said before the 
customer can choose one of the variants, if simple discounts and coupons shouldn't be used together due to business requirements.

5. Checkout the items. Usually only very small website and online shops use the logic of checking out the items from cart
and immediately go to the payment form. However, order history is vital for real e-commerce shops. It is not only necessary
for the accounting, but also is convenient for the user. Having the history of all his order enables the business:
- To get more data and recommend particular items in the future.
- To have a metrics and quickly response to the new requirements.
- To foresee the future markets.