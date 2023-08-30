# Perceptrón Multicapa

## Historia

Podemos trazar el nacimiento de las redes neuronales hasta 1958, con la introducción del Perceptrón<sup><a href="perceptron_multicapa">1</a></sup> por Frank Rosenblatt. Fue el primer algoritmo creado para reproducir la neurona biológica. Conceptualmente, el perceptrón más simple en el que se puede pensar consiste de una sola neurona: Cuando se le expone a un estímulo determinado, provee una respuesta binaria, tal como haría una neurona biológica.

![Neurona](https://i.imgur.com/uUZ5vyF.jpg)

Este modelo difiere en gran medida de la red neuronal que contiene miles de millones de neuronas en un cerebro real. Poco después de su concepción, muchos investigadores mostraron los problemas que tiene el Perceptrón; de hecho, fue rápidamente demostrado que los perceptrones no podían ser entrenados para reconocer muchos tipos de patrones de entrada. Para obtener una red más potente, era necesario disponer de múltiples niveles de unidades y crear un perceptrón multicapa, con más neuronas intermedias usadas para resolver subproblemas linealmente separables<sup><a href="linealmente_separables">2</a></sup>, cuyas salidas se combinaban en el nivel final para proporcionar una respuesta concreta al problema de entrada original. A pesar de que el Perceptrón era un clasificador binario simple pero severamente limitado, introdujo una gran innovación: la idea de simular la unidad computacional básica de un sistema biológico complejo que existe en la naturaleza.

## Teoría

Fundamentalmente, una red neuronal no es más que un muy buen aproximador funcionales; es decir, a una red neuronal se le proporciona un vector de entrada, la red realiza una serie de operaciones, y finalmente produce un vector de salida. Para entrenar una red neuronal con la finalidad de estimar una función desconocida, el proceso es muy simple: Se debe conseguir un set de entrenamiento (una colección de datos), de los cuales la red aprenderá y generalizará para hacer inferencias futuras. En un perceptrón multicapa, los datos se propagan a través de la red capa por capa hasta que llegan a la capa final. Las activaciones de la capa final son las predicciones que la red realmente hace.

### El perceptrón

Como se mencionó anteriormente, el esquema de perceptrón simple implementa una sola neurona. La forma más fácil de implementar este clasificador simple es establecer una función de umbral, insertarla en la neurona, combinar los valores (eventualmente usando diferentes pesos para cada uno de ellos) que describen el estímulo en un solo valor, proporcionar este valor a la neurona y ver qué devuelve en salida. Este esquema muestra cómo funciona:

![Perceptrón simple](https://i.imgur.com/uu0iNCC.png)

### La métrica

¿Por qué _pesos_? Conceptualmente, entrenar se refiere al "proceso de aprender las habilidades necesarias para llevar a cabo una actividad en particular". ¿Pero cómo sabemos si estamos mejorando, o si estamos aprendiendo las habilidades necesarias? Se necesita una métrica de qué tan bien o mal se está haciendo. En las redes neuronales también suele haber generalmente una métrica llamada función de costo. Supongamos que queremos cambiar cierto peso _w<sub>i</sub>_ en la red. A grosso modo, la función costo revisa la función inferida por la red y la utiliza para estimar valores para los datos en el set de entrenamiento. La diferencia entre las salidas de la red y los datos del set de entrenamiento son los valores principales para la función costo. Cuando se entrena la red, el objetivo es que el valor de esta función costo sea tan bajo como sea posible. El más básico de los algoritmos de entrenamiento es el gradiente descendente. Supongamos que podemos calcular un error _E_ basados en la variación del peso _w<sub>i</sub>_: por lo tanto, podemos dibujar el gráfico en un gráfico como el de la figura:

![Figura pesos](https://i.imgur.com/W5zQiw7.png)

De esta manera, si calculamos la derivada de esta función, podemos entender cómo la variación en el peso hace una contribución positiva o negativa al error. En la práctica, sin importar el valor derivado, podemos usar una función de corrección de peso que disminuya el peso involucrado de la cantidad derivada (modulada por la tasa de aprendizaje). A pesar de que es bastante imposible, para cualquier red o función de costo, ser verdaderamente convexa, el gradiente descendente sigue las derivadas calculadas para cada unidad de neurona para esencialmente "rodar" por la pendiente hasta que encuentre su camino hacia el centro, lo más cerca posible del mínimo global. Antes de continuar, demos un paso atrás.

### Los problemas linealmente separables

El problema del perceptrón binario hecho con una sola neurona es su incapacidad para procesar problemas no separables linealmente: Este tipo de problemas son aquellos en los que es imposible definir un hiperplano capaz de separar, en el espacio vectorial de las entradas, aquellas que requieren una salida positiva de aquellas que requieren de una salida negativa. Un ejemplo de tres puntos no colineales que pertenecen a dos clases diferentes ('+' y '-') son siempre linealmente separables en dos dimensiones. Esto se ilustra en los primeros tres ejemplos de la siguiente figura:

![Signos](https://i.imgur.com/jFC0NZR.png)

Sin embargo, no todos los conjuntos de cuatro puntos, no tres colineales, son linealmente separables en dos dimensiones. La cuarta imagen necesitaría dos líneas rectas y, por lo tanto, no es linealmente separable. Esta es la razón principal por la que los científicos comienzan a trabajar con multicapas desde el principio.

## Notas

<span id="perceptron_multicapa">1</span>: F. Rosenblatt. The perceptron: A probabilistic model for information storage and organization in the brain. Psychological Review, pages 65–386, 1958. (cit. a p. 5).

<span id="linealmente_separables">2</span>: Esta condición describe la situación en la cual existe un hiperplano que puede separar, en el espacio vectorial de las entradas, aquellas que requieren una salida positiva de aquellas que requieren una salida negativa.
