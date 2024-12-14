package locales

import "decred.org/dcrdex/client/intl"

var PtBr = map[string]*intl.Translation{
	"Language":                       {T: "pt-BR"},
	"Markets":                        {T: "Mercados"},
	"Wallets":                        {T: "Carteiras"},
	"Notifications":                  {T: "Notificações"},
	"Recent Activity":                {T: "Atividade Recente"},
	"Sign Out":                       {T: "Sair"},
	"Order History":                  {T: "Histórico de Pedidos"},
	"load from file":                 {T: "carregar do arquivo"},
	"loaded from file":               {T: "carregado do arquivo"},
	"defaults":                       {T: "padrões"},
	"Wallet Password":                {T: "Senha da Carteira"},
	"w_password_helper":              {T: "Este é a senha que você configurou com o software de sua carteira."},
	"w_password_tooltip":             {T: "Deixar senha vazia caso não haja senha necessária para sua carteira."},
	"App Password":                   {T: "Senha do App"},
	"Add":                            {T: "Adicionar"},
	"Unlock":                         {T: "Destrancar"},
	"Wallet":                         {T: "Carteira"},
	"app_password_reminder":          {T: "Sua senha do app é sempre necessária quando performando operações sensíveis da carteira"},
	"DEX Address":                    {T: "Endereço DEX"},
	"TLS Certificate":                {T: "Certificado TLS"},
	"remove":                         {T: "remover"},
	"add a file":                     {T: "adicionar um arquivo"},
	"Submit":                         {T: "Enviar"},
	"Confirm Registration":           {T: "Confirma Registro"},
	"app_pw_reg":                     {T: "Informe sua senha do app para confirmar seu registro na DEX."},
	"reg_confirm_submit":             {T: `Quando vc enviar esse formulário, <span id="feeDisplay"></span> DCR será gasto de sua carteira decred para pagar a taxa de registro.`},
	"provided_markets":               {T: "Essa DEX provê os seguintes mercados:"},
	"accepted_fee_assets":            {T: "This DEX accepts the following fees:"},
	"base_header":                    {T: "Base"},
	"quote_header":                   {T: "Quote"},
	"lot_size_headsup":               {T: `Todas as trocas são múltiplas do tamanho do lote.`},
	"Password":                       {T: "Senha"},
	"Register":                       {T: "Registrar"},
	"Authorize Export":               {T: "Autorizar exportação"},
	"export_app_pw_msg":              {T: "Informe a senha para confirmar exportação de conta"},
	"Disable Account":                {T: "Desativar Conta"},
	"disable_dex_server":             {T: "Este servidor DEX pode ser reativado a qualquer momento no futuro, adicionando-o novamente."},
	"Authorize Import":               {T: "Autorizar Importação"},
	"app_pw_import_msg":              {T: "Informe sua senha do app para confirmar importação da conta"},
	"Account File":                   {T: "Arquivo da Conta"},
	"Change Application Password":    {T: "Trocar Senha do App"},
	"Current Password":               {T: "Senha Atual"},
	"New Password":                   {T: "Nova Senha"},
	"Confirm New Password":           {T: "Confirmar Nova Senha"},
	"cancel_no_pw":                   {T: "Enviar ordem de cancelamento para o restante."},
	"cancel_remain":                  {T: "A quantidade restante pode ser alterada antes do pedido de cancelamento ser coincididos."},
	"Log In":                         {T: "Logar"},
	"epoch":                          {T: "epoque"},
	"price":                          {T: "preço"},
	"volume":                         {T: "volume"},
	"buys":                           {T: "compras"},
	"sells":                          {T: "vendas"},
	"Buy Orders":                     {T: "Pedidos de Compras"},
	"Quantity":                       {T: "Quantidade"},
	"Rate":                           {T: "Câmbio"},
	"Limit Order":                    {T: "Ordem Limite"},
	"Market Order":                   {T: "Ordem de Mercado"},
	"reg_status_msg":                 {T: `Para poder trocar em <span id="regStatusDex" class="text-break"></span>, o pagamento da taxa de registro é necessário <span id="confReq"></span> confirmações.`},
	"Buy":                            {T: "Comprar"},
	"Sell":                           {T: "Vender"},
	"lot_size":                       {T: "Tamanho do Lote"},
	"Rate Step":                      {T: "Passo de Câmbio"},
	"Max":                            {T: "Máximo"},
	"lot":                            {T: "lote"},
	"min trade is about":             {T: "troca mínima é sobre"},
	"immediate_explanation":          {T: "Se o pedido não preencher completamente durante o próximo ciclo de encontros, qualquer quantia restante não será reservada ou combinada novamente nos próximos ciclos."}, // revisar
	"Immediate or cancel":            {T: "Imediato ou cancelar"},
	"Balances":                       {T: "Balanços"},
	"outdated_tooltip":               {T: "Balanço pode está desatualizado. Conecte-se a carteira para atualizar."},
	"available":                      {T: "disponível"},
	"connect_refresh_tooltip":        {T: "Clique para conectar e atualizar"},
	"add_a_wallet":                   {T: `Adicionar uma carteira <span data-tmpl="addWalletSymbol"></span> `},
	"locked":                         {T: "trancado"},
	"immature":                       {T: "imaturo"},
	"Sell Orders":                    {T: "Pedido de venda"},
	"Your Orders":                    {T: "Seus Pedidos"},
	"Type":                           {T: "Tipo"},
	"Side":                           {T: "Lado"},
	"Age":                            {T: "Idade"},
	"Filled":                         {T: "Preenchido"},
	"Settled":                        {T: "Assentado"},
	"Status":                         {T: "Status"},
	"view order history":             {T: "ver histórico de pedidos"},
	"cancel_order":                   {T: "cancelar pedido"},
	"order details":                  {T: "detalhes do pedido"},
	"verify_order":                   {T: `Verificar<span id="vSideHeader"></span> Pedido`},
	"You are submitting an order to": {T: "Você está enviando um pedido para"},
	"at a rate of":                   {T: "Na taxa de"},
	"for a total of":                 {T: "Por um total de"},
	"verify_market":                  {T: "Está é uma ordem de mercado e combinará com o(s) melhor(es) pedidos no livro de ofertas. Baseado no atual valor médio de mercado, você receberá"}, //revisar
	"auth_order_app_pw":              {T: "Autorizar este pedido com a senha do app."},
	"lots":                           {T: "lotes"},
	"provied_markets":                {T: "Essa DEX provê os seguintes mercados:"},
	"order_disclaimer": {T: `<span class="red">IMPORTANTE</span>: Trocas levam tempo para serem concluídas, e vc não pode desligar o cliente e software DEX,
		ou o <span data-quote-ticker></span> ou <span data-base-ticker></span> blockchain e/ou software da carteira, até os pedidos serem completamente concluídos.
		A troca pode completar em alguns minutos ou levar até mesmo horas.`}, //revisar
	"Order":                      {T: "Ordem"},
	"see all orders":             {T: "ver todas as ordens"},
	"Exchange":                   {T: "Casa de Câmbio"},
	"Market":                     {T: "Mercado"},
	"Offering":                   {T: "Oferecendo"},
	"Asking":                     {T: "Pedindo"},
	"Fees":                       {T: "Taxas"},
	"order_fees_tooltip":         {T: "Taxas de transações da blockchain, normalmente coletada por mineradores. Decred DEX não coleta taxas de trocas."},
	"Matches":                    {T: "Combinações"},
	"Match ID":                   {T: "ID de Combinação"},
	"Time":                       {T: "Tempo"},
	"ago":                        {T: "atrás"},
	"Cancellation":               {T: "Cancelamento"},
	"Order Portion":              {T: "Porção do pedido"},
	"you":                        {T: "você"},
	"them":                       {T: "Eles"},
	"Redemption":                 {T: "Rendenção"},
	"Refund":                     {T: "Reembolso"},
	"Funding Coins":              {T: "Moedas de Financiamento"},
	"Exchanges":                  {T: "Casa de câmbios"}, //revisar
	"apply":                      {T: "aplicar"},
	"Assets":                     {T: "Ativos"},
	"Trade":                      {T: "Troca"},
	"Set App Password":           {T: "Definir senha de aplicativo"},
	"reg_set_app_pw_msg":         {T: "Definir senha de aplicativo. Esta senha protegerá sua conta DEX e chaves e carteiras conectadas."},
	"Password Again":             {T: "Senha Novamente"},
	"Add a DEX":                  {T: "Adicionar uma DEX"},
	"reg_ssl_needed":             {T: "Parece que não temos um certificado SSL para esta DEX. Adicione o certificado do servidor para podermos continuar."},
	"Dark Mode":                  {T: "Modo Dark"},
	"Show pop-up notifications":  {T: "Mostrar notificações de pop-up"},
	"Account ID":                 {T: "ID da Conta"},
	"Export Account":             {T: "Exportar Conta"},
	"simultaneous_servers_msg":   {T: "O cliente da DEX suporta simultâneos números de servidores DEX."},
	"Change App Password":        {T: "Trocar Senha do aplicativo"},
	"Version":                    {T: "Versão"},
	"Build ID":                   {T: "ID da Build"},
	"Connect":                    {T: "Conectar"},
	"Withdraw":                   {T: "Retirar"},
	"Deposit":                    {T: "Depositar"},
	"Lock":                       {T: "Trancar"},
	"New Deposit Address":        {T: "Novo endereço de depósito"},
	"Address":                    {T: "Endereço"},
	"Amount":                     {T: "Quantia"},
	"Reconfigure":                {T: "Reconfigurar"},
	"pw_change_instructions":     {T: "Trocando a senha abaixo não troca sua senha da sua carteira. Use este formulário para atualizar o cliente DEX depois de ter alterado a senha da carteira pela aplicação da carteira diretamente."},
	"New Wallet Password":        {T: "Nova senha da carteira"},
	"pw_change_warn":             {T: "Nota: Trocando para uma carteira diferente enquanto possui trocas ativas ou pedidos nos livros pode causar fundos a serem perdidos."},
	"Show more options":          {T: "Mostrar mais opções"},
	"seed_implore_msg":           {T: "Você deve ser cuidadoso. Escreva sua semente e salve uma cópia. Caso você perca acesso a essa maquina ou algum outra problema ocorra, você poderá usar sua semente recupear acesso a sua conta DEX e carteiras regitrada. Algumas contas antigas não podem ser recuperadas, e apesar de novas ou não, é uma boa prática salvar backup das contas de forma separada da semente."},
	"View Application Seed":      {T: "Ver semente da aplicação"},
	"Remember my password":       {T: "Lembrar senha"},
	"pw_for_seed":                {T: "Informar sua senha do aplicativo para mostrar sua seed. Tenha certeza que mais ninguém pode ver sua tela."},
	"Asset":                      {T: "Ativo"},
	"Balance":                    {T: "Balanço"},
	"Actions":                    {T: "Ações"},
	"Restoration Seed":           {T: "Restaurar Semente"},
	"Restore from seed":          {T: "Restaurar da Semente"},
	"Import Account":             {T: "Importar Conta"},
	"no_wallet":                  {T: "sem carteira"},
	"create_a_x_wallet":          {T: "Criar uma Carteira <span data-asset-name=1></span>"},
	"dont_share":                 {T: "Não compartilhe e não perca sua seed."},
	"Show Me":                    {T: "Mostre me"},
	"Wallet Settings":            {T: "Configurações da Carteira"},
	"add_a_x_wallet":             {T: `Adicionar uma carteira <img data-tmpl="assetLogo" class="small-icon mx-1"> <span data-tmpl="assetName"></span>`},
	"ready":                      {T: "destrancado"},
	"off":                        {T: "desligado"},
	"Export Trades":              {T: "Exportar Trocas"},
	"change the wallet type":     {T: "trocar tipo de carteira"},
	"confirmations":              {T: "confirmações"},
	"pick a different asset":     {T: "Escolher ativo diferente"},
	"Create":                     {T: "Criar"},
	"1 Sync the Blockchain":      {T: "1: Sincronizar a Blockchain"},
	"Progress":                   {T: "Progresso"},
	"remaining":                  {T: "Faltando"},
	"Your Deposit Address":       {T: "Seu Endereço de Depósito"},
	"add a different server":     {T: "Adicionar um servidor diferente"},
	"Add a custom server":        {T: "Adicionar um servidor personalizado"},
	"plus tx fees":               {T: "+ tx fees"},
	"Export Seed":                {T: "Exportar Seed"},
	"Total":                      {T: "Total"},
	"Trading":                    {T: "Trocando"},
	"Receiving Approximately":    {T: "Recebendo aproximadamente"},
	"Fee Projection":             {T: "Projeção de Taxa"},
	"details":                    {T: "detalhes"},
	"to":                         {T: "para"},
	"Options":                    {T: "Opções"},
	"fee_projection_tooltip":     {T: "Se as condições da rede não mudarem antes que seu pedido corresponda, O total de taxas realizadas deve estar dentro dessa faixa."},
	"unlock_for_details":         {T: "Desbloqueie suas carteiras para recuperar detalhes do pedido e opções adicionais."},
	"estimate_unavailable":       {T: "Orçamentos e opções de pedidos indisponíveis"},
	"Fee Details":                {T: "Detalhes da taxa"},
	"estimate_market_conditions": {T: "As estimativas de melhor e pior caso são baseadas nas condições atuais da rede e podem mudar quando o pedido corresponder."},
	"Best Case Fees":             {T: "Cenário de melhor taxa"},
	"best_case_conditions":       {T: "O melhor cenário para taxa ocorre quando todo pedido é correspondido em uma única combinação."},
	"Swap":                       {T: "Troca"},
	"Redeem":                     {T: "Resgatar"},
	"Worst Case Fees":            {T: "Pior cenário para taxas"},
	"worst_case_conditions":      {T: "O pior caso pode ocorrer se a ordem corresponder em um lote de cada vez ao longo de muitos epoques."},
	"Maximum Possible Swap Fees": {T: "Taxas Máximas de Troca Possíveis"},
	"max_fee_conditions":         {T: "Este é o máximo que você pagaria em taxas em sua troca. As taxas são normalmente avaliadas em uma fração dessa taxa. O máximo não está sujeito a alterações uma vez que seu pedido é feito."},
}
