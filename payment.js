document.addEventListener('DOMContentLoaded',async()=>{

  const Publishable_key=  await fetch("/config").then(r=>r.json());
  const stripe=Stripe(Publishable_key);


    //fetch the payment intnets client secret

    const{clientSecret}=await fetch("payment",{
        method: 'POST',
        headers: {
            "Content-Type": "application/json",
        },

    }).then(r=>r.json())


    //Mount the Elements

    const elements=stripe.elements({clientSecret})
    const paymentElement=elements.create('payment')
    paymentElement.mount('#payment-element')



})