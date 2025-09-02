# Regras Oficiais — Jogo de Cartas: Pedra • Papel • Tesoura

## 1. Componentes
- Baralho de cartas do jogo contendo três tipos: **Pedra**, **Papel**, **Tesoura**.  
- Cada carta possui um **nível de melhoria (estrelas)** inteiro: 0 (sem estrelas) até 5 estrelas.  
- Cada jogador tem um **nível de jogador** (`level`) e um valor de **HP máximo** ligado a esse nível. O nível inicial padrão é 1.  
- Loja externa: permite obter cartas / realizar fusões conforme o sistema de melhoria (detalhado na seção 6).

---

## 2. Objetivo
Eliminar o HP do oponente. O jogador que ficar com HP ≤ 0 perde a partida.

---

## 3. Preparação
- Cada jogador começa com HP igual ao HP máximo definido por seu `level`.  
- **HP máximo por nível**:  
HP_max(l) = 500 × 2^(l - 1)

- Nível 1 → 500 HP  
- Nível 2 → 1000 HP  
- Nível 3 → 2000 HP  
- ...

---

## 4. Ordem e sequência de turno
1. Em cada turno, **simultaneamente** ambos os jogadores escolhem **uma** carta da mão (Pedra, Papel ou Tesoura) e a revelam.  
2. A resolução do turno segue as regras descritas na seção 5.  
3. Não há iniciativa: ambas as cartas se confrontam e o dano (se houver) é aplicado conforme o resultado do confronto.

---

## 5. Resolução de confronto (por turnos)

### 5.1. Resultado por tipo de carta (RPS clássico)
- Se as cartas forem **diferentes**, aplica-se a regra padrão:
- Pedra < Papel < Tesoura < Pedra  
- O **vencedor pelo tipo** inflige **100% do seu dano de ataque emitido**.  
- O perdedor sofre a redução de HP equivalente.

### 5.2. Empate por tipo (mesmo tipo jogado)
- Se ambos jogarem o **mesmo tipo**, compara-se o número de **estrelas**:
- Carta com **mais estrelas** vence.  
- Se diferentes, calcula-se o dano por **diferença entre ataque e defesa** (ver seção 7).  
- Se iguais, **o turno é pulado**.

### 5.3. Regras adicionais
- **Dano final mínimo = 0** (nunca negativo).  
- Se `dano_final ≥ HP atual do defensor`, o defensor vai a 0 HP e perde imediatamente.

---

## 6. Sistema de melhoria de cartas (fusões e estrelas)

### 6.1. Upgrades incrementais
- Estrelas vão de 0 → 1 → 2 → 3 → 4 → 5.  
- Para subir uma estrela (`s → s+1`):
- Usar **(s+1)** cartas de estrela `s`, **OU**  
- Usar **1** carta de estrela `(s+1)`.

**Exemplos:**
- 0 → 1: 1 carta de 0★ (ou 1 carta de 1★)  
- 1 → 2: 2 cartas de 1★ (ou 1 carta de 2★)  
- 2 → 3: 3 cartas de 2★ (ou 1 carta de 3★)  
- 3 → 4: 4 cartas de 3★ (ou 1 carta de 4★)  
- 4 → 5: 5 cartas de 4★ (ou 1 carta de 5★)  

### 6.2. Saltos múltiplos
- Permitidos apenas se:
- O jogador possuir **1 carta da estrela alvo**, **OU**  
- Fornecer a equivalência em cartas menores conforme a cadeia.  

**Exemplo (2 → 4):**
- 1 carta de 4★  
- 4 cartas de 3★  
- 12 cartas de 2★  

### 6.3. Limites
- Máximo de estrelas: **5**.

---

## 7. Fórmulas de dano e defesa

### 7.1. Dano emitido
a(l, s, h) = (0.1 × h + 20 × s) × (1.005^l)

- `l` = nível do atacante  
- `s` = estrelas da carta atacante  
- `h` = HP máximo do defensor  

### 7.2. Defesa (somente em desempate por estrelas)
defesa(l, s) = 20 × s × (1.05^l)

- `l` = nível do defensor  
- `s` = estrelas da carta defensora  

### 7.3. Dano final
- Vitória por **tipo**:  
dano_final = a
- Vitória por **estrelas**:  
dano_final = max(0, a - defesa)

### 7.4. Observações
- Carta 0★: dano base = `10% do HP do oponente`.  
- Carta 5★: dano base = `10% do HP + 100`.  

---

## 8. Experiência e progressão de nível

### 8.1. Ganho de EXP
- **Vencedor**: `10% × HP_max do inimigo`.  
- **Perdedor**: `2.5% × HP_max do inimigo`.

### 8.2. HP por nível
HP_max(l) = 500 × 2^(l - 1)

### 8.3. EXP necessária
- 1 → 2 = 100  
- 2 → 3 = 400  
- 3 → 4 = 1000  
- 4 → 5 = 1900  
- 5 → 6 = 3300  

---

## 9. Condições de término
- Perde quem tiver HP = 0.  
- Se ambos chegam a 0 no mesmo turno → **empate** (salvo regra adicional definida).

---

## 10. Resumo de turno
1. Ambos escolhem e revelam carta.  
2. Se tipos diferentes → aplica RPS → vencedor causa dano `a`.  
3. Se mesmo tipo → compara estrelas:  
   - mais estrelas vence → aplica `max(0, a - defesa)`  
   - mesmas estrelas → turno pulado.  
4. Atualizar HP.  
5. Se alguém chegou a 0 HP → partida termina.  
6. No fim da partida, aplicar EXP conforme resultado.
