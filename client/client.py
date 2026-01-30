# client.py

import grpc
import order_pb2
import order_pb2_grpc

def enviar_pedido(stub, nome_teste, customer_id, itens):
    """
    Função auxiliar para enviar o pedido e imprimir o resultado formatado.
    """
    print(f"🔵 EXECUTANDO: {nome_teste}")
    
    # Monta a requisição
    request = order_pb2.CreateOrderRequest(
        costumer_id=customer_id,
        order_items=itens
    )

    try:
        # Tenta enviar para o microsserviço Order
        response = stub.Create(request)
        print(f"✅ SUCESSO! Pedido criado com ID: {response.order_id}")
        print("   -> Status deve ser 'Shipped' no banco de dados.")
    
    except grpc.RpcError as e:
        # Captura o erro retornado pelo servidor
        print(f"❌ ERRO RECEBIDO (Status gRPC: {e.code()})")
        print(f"   Mensagem: {e.details()}")
        
        error_msg = e.details() if e.details() else ""
        
        if "exceed 50" in error_msg:
             print("   -> CORRETO - Bloqueio de quantidade funcionou.")
        
        elif "Payment over 1000" in error_msg:
             print("   -> CORRETO -Bloqueio de pagamento funcionou.")
             print("   -> EXPECTATIVA: Status no banco deve ser 'Canceled'")
        
        elif "product not found" in error_msg:
             print("   -> CORRETO -Validação de Estoque funcionou (Requisito 1.2).")
        
        elif "connectex" in error_msg or "unavailable" in error_msg.lower():
             print("   -> SINAL DE ALERTA: Parece que um dos microsserviços está desligado.")
        else:
             print("   -> Erro não esperado.")
    
    print("-" * 50 + "\n")


def run():
    # Conectar ao servidor gRPC (ajuste o host/porta conforme necessário)
    # Conectar ao servidor gRPC na porta 3000 (Order Service)
    print("🔌 Conectando ao servidor gRPC...")
    channel = grpc.insecure_channel('localhost:3000')
    stub = order_pb2_grpc.OrderStub(channel)
    print("-" * 50 + "\n")

    item_valido = order_pb2.OrderItem(
        product_code="CANETA",
        unit_price=10.0,
        quantity=5
    )
    enviar_pedido(stub, "Teste: Pedido Válido", 101, [item_valido])


    # Erro de Pagamento - Preço > 1000. - PAYMENT
    item_caro = order_pb2.OrderItem(
        product_code="NOTEBOOK",
        unit_price=1500.0, #ultrapassando limite de preço
        quantity=1
    )
    enviar_pedido(stub, "Teste 3: Preço Alto (> R$ 1000)", 103, [item_caro])


if __name__ == '__main__':
    run()